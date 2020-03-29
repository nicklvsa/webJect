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

//Decompress - decompresses and walks a tree
func Decompress(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else if f.FileInfo().Mode()&os.ModeSymlink != 0 {
			buffer := make([]byte, f.FileInfo().Size())
			size, err := rc.Read(buffer)
			if err != nil {
				return err
			}

			target := string(buffer[:size])

			err = os.Symlink(target, path)
			if err != nil {
				return err
			}
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err = f.Close(); err != nil {
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

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

//RemoveWalkDirs - removes directory recursively
func RemoveWalkDirs(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
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
