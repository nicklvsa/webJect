package shared

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/xid"
)

var (
	logger *Logger
)

//BuildResponse - builds and returns json encoded endpoint response
func BuildResponse(errCode int, contentMsg, errMsg interface{}, w http.ResponseWriter) []byte {

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin")
	w.Header().Set("Content-Type", "application/json")

	apiResponse := &APIResponse{
		ErrorCode:      errCode,
		ErrorMessage:   errMsg,
		ContentMessage: contentMsg,
	}

	response, err := json.Marshal(apiResponse)
	if err != nil {
		logger.Err("Could not build the api response!")
	}

	logger.Info(response)
	return response
}

//IsWindows - returns if the operating system is windows or not
func IsWindows() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

//IsMacOS - returns if the operating system is macos or not
func IsMacOS() bool {
	return strings.Contains(runtime.GOOS, "mac")
}

//Decompress - unzips a zip file
func Decompress(src, dest string) error {
	read, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	defer func() {
		if err := read.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	extractAndWrite := func(file *zip.File) error {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), file.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}

			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range read.File {
		err := extractAndWrite(f)
		if err != nil {
			return err
		}
	}

	return nil
}

//CleanDecompression - removes the zip files after extraction
func CleanDecompression(name, path string) error {
	err := os.Remove(path + string(os.PathSeparator) + name)
	if err != nil {
		return err
	}
	return nil
}

//GenID - returns a unique id
func GenID() string {
	return xid.New().String()
}

//Info - logs an info message
func (log *Logger) Info(content interface{}) {
	fmt.Println(fmt.Sprintf("[Info] - %+v", content))
}

//Warn - logs a warn message
func (log *Logger) Warn(content interface{}) {
	fmt.Println(fmt.Sprintf("[WARN] - %+v", content))
}

//Err - logs an error message
func (log *Logger) Err(content interface{}) {
	fmt.Println(fmt.Sprintf("[ERROR] - %+v", content))
}
