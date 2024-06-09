package sc

import "path"

const location string = "/mnt/fastbulk/Games/star-citizen/drive_c/Program Files/Roberts Space Industries/StarCitizen/LIVE/Data.p4k"

type Version string

const (
	EPTU Version = "EPTU"
	PTU  Version = "PTU"
	Live Version = "LIVE"
)

func ClientPath(prefix string, version Version) string {
	return path.Join(BasePath(prefix), "StarCitizen", string(version))
}

func UserPath(prefix string, version Version) string {
	return path.Join(ClientPath(prefix, version), "USER", "Client", "0")
}

func Profiles(prefix string, version Version) string {
	return path.Join(UserPath(prefix, version), "Profiles", "default")
}
func P4K(prefix string, version Version) string {
	return path.Join(ClientPath(prefix, version), "Data.p4k")
}

func ActionMaps(prefix string, version Version) string {
	return path.Join(Profiles(prefix, version), "actionmaps.xml")
}
