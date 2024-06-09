package sc

func BasePath(prefix string) string {
	if prefix == "" {
		return path.Join("C:", "Program Files", "Roberts Space Industries")
	}
	return path.Join(prefix, "Roberts Space Industries")
}
