package devops

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func NormalizeLocalPath(userInputPath string) (string, error) {
	pathOfInterest := userInputPath

	if strings.Contains(userInputPath, "~") && userInputPath[0] == '~' {
		homeDirectoryPath, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve ~ to user home directory: %s", err)
		}
		pathOfInterest = path.Join(homeDirectoryPath, strings.Replace(userInputPath, "~", "", 1))
		if strings.Contains(pathOfInterest, "~") {
			return "", fmt.Errorf("refusing to replace multiple ~ references")
		}
	}

	if !path.IsAbs(pathOfInterest) {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to resolve . to current directory: %s", err)
		}
		pathOfInterest = path.Join(currentWorkingDirectory, pathOfInterest)
	}

	return pathOfInterest, nil
}
