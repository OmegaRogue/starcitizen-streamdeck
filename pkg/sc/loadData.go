package sc

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"starcitizen-streamdeck/internal/util"
	cryxml2 "starcitizen-streamdeck/pkg/cryxml"
	"starcitizen-streamdeck/pkg/p4k"
)

func LoadData(prefix string, version Version) Data {
	defaultProfile, _ := os.OpenFile("defaultProfile.xml", os.O_RDWR|os.O_CREATE, 0666)
	defer util.DiscardErrorOnly(defaultProfile.Sync())
	defer util.DiscardErrorOnly(defaultProfile.Close())

	global, _ := os.OpenFile("global.json", os.O_RDWR|os.O_CREATE, 0666)
	defer util.DiscardErrorOnly(global.Sync())
	defer util.DiscardErrorOnly(global.Close())

	p4kPath := P4K(prefix, version)

	p4kFi, _ := os.Stat(p4kPath)
	pd := p4k.NewDirectory()
	var cryXMLData string
	if fi, _ := defaultProfile.Stat(); fi.Size() <= 10 || p4kFi.ModTime().Before(fi.ModTime()) {
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
		tree.BuildXml(root)
		cryXMLData = fmt.Sprint(tree)

		if _, err := io.WriteString(defaultProfile, cryXMLData); err != nil {
			log.Fatal().Err(err).Msg("error writing defaultProfile")
		}
	} else {
		data, _ := io.ReadAll(defaultProfile)
		cryXMLData = string(data)
	}

	locals := map[string]map[string]string{}
	if fi, _ := global.Stat(); fi.Size() <= 10 || p4kFi.ModTime().Before(fi.ModTime()) {
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

	//log.Info().Err(err).Str("filename", file.Filename).Msg("\n" + retVal)

	profile, _ := FromXml(cryXMLData)

	return Data{
		profile, LoadActionmap(prefix, version),
		locals,
	}
}

func LoadActionmap(prefix string, version Version) *ActionMapActionMaps {
	data, _ := os.ReadFile(ActionMaps(prefix, version))
	//fmt.Println(string(data))

	var out ActionMapActionMaps
	if err := xml.Unmarshal(data, &out); err != nil {
		log.Fatal().Err(err).Msg("error unmarshalling ActionMaps")
	}
	out.Prepare()
	return &out
}
