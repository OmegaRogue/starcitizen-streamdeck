//go:generate go-enum -f=$GOFILE --marshal --nocase -t ../../files/zerolog.gotmpl
package sc

import "path"

// Version
/*
ENUM(
EPTU
PTU
LIVE
)
*/
type Version string

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
