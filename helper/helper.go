package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func OpenFile(path string) (os.File, error) {
	opened, openingFileError := os.Open(path)

	if openingFileError != nil {
		return os.File{}, openingFileError
	}

	return *opened, nil
}

func CreateFile(fileName string, path string) (os.File, error) {
	file, fileCreateError := os.Create(fileName)

	if fileCreateError != nil {
		return os.File{}, fileCreateError
	}

	defer file.Close()

	return *file, nil
}

func WriteFile(file *os.File, text string) (bool, error) {
	_, writtingError := file.Write([]byte(text))
	if writtingError != nil {
		return false, writtingError
	}
	defer file.Close()

	return true, nil
}

func ReadFile(filepath string) (string, error) {
	data, readingError := os.ReadFile(filepath)

	if readingError != nil {
		return "", readingError
	}

	return string(data), nil
}

func DeleteFile(path string) bool {
	return os.Remove(path) == nil
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return (err == nil && !os.IsNotExist(err))
}

func CreateDirectory(path string, mode os.FileMode) bool {
	if Exists(path) {
		err := os.Mkdir(path, mode)

		return err == nil
	}

	return false
}

func HashString(data string) (string, error) {
	hashData, hashError := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)

	if hashError != nil {
		return "", hashError
	}

	return string(hashData), nil
}

func CheckPassword(passFromDatabase, passFromInput string) bool {
	passCheckError := bcrypt.CompareHashAndPassword([]byte(passFromDatabase), []byte(passFromInput))

	return passCheckError == nil
}

func BaseURL(r *http.Request) string {
	// Determine the scheme (http or https)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Construct the base URL
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func ResponseJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func FormateDateTime(realDate, formta string) (string, error) {
	strTiem, timeError := time.Parse(time.RFC3339, realDate)

	if timeError != nil {
		panic(timeError)
	}

	return strTiem.Format(formta), nil
}

func HasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)

		if filepath.Ext(file) == ext {
			return true
		}
	}

	return false
}
