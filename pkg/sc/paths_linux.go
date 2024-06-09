package sc

import "path"

func BasePath(prefix string) string {
	return path.Join(prefix, "drive_c", "Program Files", "Roberts Space Industries")
}
