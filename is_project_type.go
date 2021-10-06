package devops

import (
	"fmt"
	"io/ioutil"
)

type ProjectType string

const (
	TypeC          ProjectType = "c"
	TypeGo         ProjectType = "go"
	TypeHaskell    ProjectType = "hs"
	TypeJava       ProjectType = "java"
	TypeJavascript ProjectType = "js"
	TypePython     ProjectType = "py"
	TypeRuby       ProjectType = "rb"
	TypeRust       ProjectType = "rs"
	TypeTypescript ProjectType = "ts"
)

var projectFileMap = map[ProjectType][]string{
	TypeC:          {"configure", "configure.ac"},
	TypeGo:         {"go.mod", "go.sum"},
	TypeHaskell:    {"stack.yaml", "stack.yaml.lock"},
	TypeJava:       {"gradlew"},
	TypeJavascript: {"package.json", "package-lock.json", "yarn.lock"},
	TypePython:     {"requirements.txt", "tox.ini"},
	TypeRuby:       {"Gemfile", "Gemfile.lock"},
	TypeRust:       {"Cargo.lock", "Cargo.toml"},
	TypeTypescript: {"tsconfig.json"},
}
var projectDirMap = map[ProjectType][]string{
	TypeGo:         {"vendors"},
	TypeJavascript: {"node_modules"},
}

func IsProjectType(pathToDirectory string, projectType ProjectType) (bool, error) {
	normalizedPath, err := NormalizeLocalPath(pathToDirectory)
	if err != nil {
		return false, fmt.Errorf("failed to normalize input path '%s': %s", pathToDirectory, err)
	}
	fileListings, err := ioutil.ReadDir(normalizedPath)
	if err != nil {
		return false, fmt.Errorf("failed to access normalized path '%s': %s", normalizedPath, err)
	}
	directories := []string{}
	files := []string{}
	for _, fileListing := range fileListings {
		if fileListing.IsDir() {
			directories = append(directories, fileListing.Name())
		} else {
			files = append(files, fileListing.Name())
		}
	}
	if containsAnyString(files, projectFileMap[projectType]) {
		return true, nil
	}
	if containsAnyString(directories, projectDirMap[projectType]) {
		return true, nil
	}
	return false, nil
}
