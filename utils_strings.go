package devops

func containsAllStrings(haystack []string, needles []string) bool {
	if haystack == nil || needles == nil || len(haystack) == 0 || len(needles) == 0 {
		return false
	}
	needleMap := map[string]bool{}
	for _, needle := range needles {
		needleMap[needle] = false
	}
	for _, item := range haystack {
		if _, ok := needleMap[item]; ok {
			needleMap[item] = true
		}
	}
	allStringsExist := true
	for _, exists := range needleMap {
		allStringsExist = allStringsExist && exists
	}
	return allStringsExist
}

func containsAnyString(haystack []string, needles []string) bool {
	if haystack == nil || needles == nil || len(haystack) == 0 || len(needles) == 0 {
		return false
	}
	needleMap := map[string]bool{}
	for _, needle := range needles {
		needleMap[needle] = true
	}
	for _, item := range haystack {
		if exists, ok := needleMap[item]; exists && ok {
			return true
		}
	}
	return false
}
