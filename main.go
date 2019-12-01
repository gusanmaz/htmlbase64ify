/* This CLI encodes source images of html and markdown files as base64 data. By using this CLI web data could be stored and shared in a more compact fashion. Referred image files and source html file should be in the same directory for CLI to be able to convert source images as base64 data. This restriction may be loosened in future editions. */
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ImgExtensionsReg stores image file extensions as regexp to match image files.
var ImgExtensionsReg = "(jpe?g)|(png)|(gif)|(bmp)|(tiff)$"

/* IsImageFile detects whether given filename bears an image format extension as suffix.
It doesn't detect whether given filename exists in the file system. */
func IsImageFile(fileName string) bool {
	imRegexp := regexp.MustCompile(ImgExtensionsReg)
	return imRegexp.MatchString(fileName)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please specify a html file name as a parameter!")
		os.Exit(1)
	}

	fileName := args[1]
	info, err := os.Stat(fileName)
	srcFileMode := info.Mode()
	if os.IsNotExist(err) || info.IsDir() {
		fmt.Println("Please specify a valid html file." + fileName + " is not a valid file name!")
		os.Exit(1)
	}

	destFile := fileName
	if len(args) >= 3 {
		destFile = args[2]
		destFileDir := filepath.Dir(destFile)
		_, err := os.Stat(destFileDir)
		if os.IsNotExist(err) || destFile == destFileDir {
			fmt.Println("Please specify a valid destination file." + fileName + " is not a valid file name!")
			os.Exit(1)
		}
	}

	fileDir := filepath.Dir(fileName)
	files, _ := ioutil.ReadDir(fileDir)
	//imgBase64 := make(map[string]string)
	htmlFileContent, _ := ioutil.ReadFile(fileName)
	htmlFileContentS := string(htmlFileContent)

	for _, file := range files {
		if IsImageFile(file.Name()) {
			fileContent, _ := ioutil.ReadFile(file.Name())
			content64 := base64.StdEncoding.EncodeToString([]byte(fileContent))
			content64 = fmt.Sprintf("%s%s", "data:image/jpeg;base64,", content64)
			//imgBase64[fileName] = content64
			htmlFileContentS = strings.ReplaceAll(htmlFileContentS, file.Name(), content64)
		}
	}

	ioutil.WriteFile(destFile, []byte(htmlFileContentS), srcFileMode)
}
