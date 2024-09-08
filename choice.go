package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ChooseFile checks supplied argument. If it is name of exisiting file with
// extensions "jpg" or "jpeg" it just returns its back. If it is directory, it
// chooses image file ("jpg" or "jpeg") from supplied directory and returns
// its name with the path. Otherwise it returns empty string and the error.
func ChooseFile(execPath, name string) (string, error) {
	// first look if it is "jpeg" "jpg" file
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	if !fi.IsDir() { // this is filename, just return it
		return name, nil
	}
	// this is the directory. Select the image file
	fis, err := os.ReadDir(name)
	if err != nil {
		return "", err
	}
	// copy into the folder names, which are files with extensions jpg/jpeg
	var jpgs []string
	for _, elem := range fis {
		fi, err = os.Stat(filepath.Join(name, elem.Name()))
		if err != nil {
			continue
		}
		if fi.IsDir() {
			continue
		}
		extension := strings.ToLower(filepath.Ext(elem.Name()))
		if extension != ".jpg" && extension != ".jpeg" {
			continue
		}
		jpgs = append(jpgs, elem.Name())
	}
	if len(jpgs) == 0 {
		return "", fmt.Errorf("no image JPG in the directory %s", name)
	}
	// get from the database
	imageName, ok := GetImageName(execPath, name, jpgs)
	if ok {
		return filepath.Join(name, imageName), nil
	}
	return "", fmt.Errorf("error")
}
