package test

import (
	"course/helper"
	"path/filepath"
	"testing"
)

var fileName string = "createfile.json"
var rootPath string = helper.RootDir()
var directoryPath string = filepath.Join(rootPath, "testingFile")
var filePath string = filepath.Join(directoryPath, fileName)

func TestCreateFile(t *testing.T) {
	helper.CreateDirectory(directoryPath, 0777)

	file, createError := helper.CreateFile(fileName, filePath)

	if createError != nil {
		t.Error(createError)
	}

	t.Log(file)
}

func TestFileOrFolderExistFunction(t *testing.T) {
	var existsBoolean bool = helper.Exists(directoryPath)

	if !existsBoolean {
		t.Fatalf("File not exists")
	}
}

func TestReadingDataFromFile(t *testing.T) {
	data, readingError := helper.ReadFile(filePath)

	if readingError != nil {
		t.Error(readingError)
	}

	t.Log(data)
}
