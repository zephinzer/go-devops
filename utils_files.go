package devops

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func recursivelyGetExtensionsCount(pathToDirectory string, ignoreList ...string) (map[string]int, error) {
	errors := []string{}
	results := map[string]int{}
	ignoredList := []regexp.Regexp{}
	fileEntries, err := ioutil.ReadDir(pathToDirectory)
	if err != nil {
		return results, fmt.Errorf("failed to list directory contents at '%s': %s", pathToDirectory, err)
	}
	for _, ignorable := range ignoreList {
		ignoredList = append(ignoredList, *regexp.MustCompile(ignorable))
	}

	for _, fileEntry := range fileEntries {
		filename := fileEntry.Name()
		fullPath := path.Join(pathToDirectory, filename)
		isIgnored := false
		for _, ignorable := range ignoredList {
			if ignorable.MatchString(filename) {
				isIgnored = true
				break
			}
		}
		if isIgnored {
			continue
		}
		if fileEntry.IsDir() {
			nextMatches, err := recursivelyGetExtensionsCount(fullPath, ignoreList...)
			if err != nil {
				errors = append(errors, err.Error())
			}
			for nextMatch, count := range nextMatches {
				results[nextMatch] += count
			}
		} else {
			ext := path.Ext(filename)
			results[ext] += 1
		}
	}
	if len(errors) > 0 {
		return results, fmt.Errorf("failed to successfully get all files with extension: ['%s']", strings.Join(errors, "', '"))
	}
	return results, nil
}
