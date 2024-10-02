package models

import (
	"course/helper"
	"encoding/json"
	"html/template"
	"path/filepath"
)

var fileName string = "createfile.json"
var rootPath string = helper.RootDir()
var directoryPath string = filepath.Join(rootPath, "testingFile")
var FilePathMail string = filepath.Join(directoryPath, fileName)

type FAQ struct {
	Question   template.HTML `json:"Question"`
	Answer     template.HTML `json:"Answer"`
	UserRating Ratings       `json:"UserRating"`
}

type Ratings struct {
	User  string `json:"User"`
	Email string `json:"Email"`
	Image string `json:"Image"`
}

func (f FAQ) Get() []FAQ {
	data, readingError := helper.ReadFile(FilePathMail)

	if readingError != nil {
		panic(readingError)
	}

	var faqs []FAQ = make([]FAQ, 0)

	jsonError := json.Unmarshal([]byte(data), &faqs)

	if jsonError != nil {
		panic(jsonError)
	}
	return faqs
}

func (f FAQ) FindSingle(id int) FAQ {
	var faqs []FAQ = f.Get()

	return faqs[id-1]
}

func (f FAQ) GetRealFilePath() string {
	return FilePathMail
}
