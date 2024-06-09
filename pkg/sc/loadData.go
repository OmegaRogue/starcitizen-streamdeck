package sc

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	cryxml2 "starcitizen-streamdeck/pkg/cryxml"
	"starcitizen-streamdeck/pkg/p4k"
)

func LoadData(prefix string, version Version) Data {
	defaultExists := false
	globalExists := false

	p4kPath := P4K(prefix, version)
	p4kFi, _ := os.Stat(p4kPath)

	defaultProfile, err := os.OpenFile("defaultProfile.xml", os.O_RDWR|os.O_CREATE, 0o666)
	if err == nil {
		fi, _ := defaultProfile.Stat()
		defaultExists = fi.Size() <= 10 || p4kFi.ModTime().Before(fi.ModTime())
	}
	global, err := os.OpenFile("global.json", os.O_RDWR|os.O_CREATE, 0o666)
	if err == nil {
		fi, _ := global.Stat()
		globalExists = fi.Size() <= 10 || p4kFi.ModTime().Before(fi.ModTime())
	}

	pd := p4k.NewDirectory()
	var cryXMLData string
	if defaultExists {
		file, err := pd.ScanDirectoryFor(p4kPath, "defaultProfile.xml")
		if err != nil {
			log.Fatal().Err(err).Msg("error scanning p4k")
		}
		fd, err := p4k.GetFile(p4kPath, file)
		if err != nil {
			log.Fatal().Err(err).Msg("error reading file")
		}
		cbr := new(cryxml2.BinReader)
		root, err := cbr.LoadFromBuffer(fd)
		if err != nil {
			log.Fatal().Err(err).Msg("error parsing cryxml")
		}
		tree := new(cryxml2.Tree)
		tree.BuildXML(root)
		cryXMLData = tree.String()

		if _, err := defaultProfile.WriteString(cryXMLData); err != nil {
			log.Fatal().Err(err).Msg("error writing defaultProfile")
		}
	} else {
		data, _ := io.ReadAll(defaultProfile)
		cryXMLData = string(data)
	}

	locals := map[string]map[string]string{}
	if globalExists {
		localFiles, err := pd.ScanDirectoryContaining(p4kPath, "english\\global.ini")
		if err != nil {
			log.Fatal().Err(err).Msg("error scanning p4k")
		}

		for _, localFile := range localFiles {
			localFileData, err := p4k.GetFile(p4kPath, localFile)
			if err != nil {
				log.Fatal().Err(err).Msg("error reading file")
			}
			name := strings.TrimPrefix(strings.TrimSuffix(localFile.Filename, "\\global.ini"), "Data\\Localization\\")
			locals[name] = lo.Associate(strings.Split(string(localFileData), "\r\n"), func(item string) (string, string) {
				key, value, _ := strings.Cut(strings.TrimSpace(item), "=")
				return key, value
			})
		}
		out, _ := json.Marshal(locals)
		if _, err := global.Write(out); err != nil {
			log.Fatal().Err(err).Msg("error writing globalization")
		}
	} else {
		data, _ := io.ReadAll(global)
		if err := json.Unmarshal(data, &locals); err != nil {
			log.Fatal().Err(err).Msg("error unmarshalling globalization")
		}
	}

	profile, _ := FromXml(cryXMLData)

	if defaultProfile != nil {

		if err := defaultProfile.Sync(); err != nil {
			log.Fatal().Err(err).Msg("error syncing defaultProfile")
		}
		if err := defaultProfile.Close(); err != nil {
			log.Fatal().Err(err).Msg("error closing defaultProfile")
		}
	}
	if global != nil {
		if err := global.Sync(); err != nil {
			log.Fatal().Err(err).Msg("error syncing global")
		}
		if err := global.Close(); err != nil {
			log.Fatal().Err(err).Msg("error closing global")
		}
	}
	return Data{
		profile, LoadActionmap(prefix, version),
		locals,
	}
}

func LoadActionmap(prefix string, version Version) *ActionMapActionMaps {
	data, _ := os.ReadFile(ActionMaps(prefix, version))

	var out ActionMapActionMaps
	if err := xml.Unmarshal(data, &out); err != nil {
		log.Fatal().Err(err).Msg("error unmarshalling ActionMaps")
	}
	out.Prepare()
	return &out
}
