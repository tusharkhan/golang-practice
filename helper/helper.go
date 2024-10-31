package helper

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"

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
