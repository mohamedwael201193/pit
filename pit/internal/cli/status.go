package cli

func StatusNeverSigns(line string) bool {
	return line != "" && !containsSecret(line)
}

func containsSecret(line string) bool {
	return RefusePrint(line) != nil
}

func LinkCopy(linked bool, err error) string {
	if err != nil {
		return "agent    query_failed"
	}
	if linked {
		return "agent    linked"
	}
	return "agent    approveAgent_required"
}

func VenueCopy(found bool, err error) string {
	if err != nil {
		return "orders   query_failed"
	}
	if found {
		return "orders   on_venue"
	}
	return "orders   not_on_venue"
}
