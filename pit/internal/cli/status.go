package cli

func StatusNeverSigns(line string) bool {
	return line != "" && !containsSecret(line)
}

func containsSecret(line string) bool {
	return RefusePrint(line) != nil
}
