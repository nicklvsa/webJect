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
	return strings.Contains(runtime.GOOS, "darwin")
}

//Decompress - unzips a zip file
func Decompress(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

//CleanDecompression - removes the zip files after extraction
func CleanDecompression(name string) error {
	err := os.Remove(name)
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
