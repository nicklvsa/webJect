package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"webject/bridge"
	"webject/shared"

	"github.com/gorilla/mux"
)

var (
	logger *shared.Logger
)

//BaseAPIHandler - base / route
func BaseAPIHandler(w http.ResponseWriter, r *http.Request) {
	response := shared.BuildResponse(0, "Hello, World!", nil, w)
	logger.Info(string(response))
	w.Write(response)
}

//TweakAPIHandler - route that handles the tweak config files
func TweakAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	action := vars["action"]

	switch action {
	case "add":

		r.Body = http.MaxBytesReader(w, r.Body, shared.MaxUpload)
		if err := r.ParseMultipartForm(shared.MaxUpload); err != nil {
			response := shared.BuildResponse(1, nil, fmt.Sprintf("File too big to handle! Error: %s", err.Error()), w)
			w.Write(response)
			return
		}

		file, _, err := r.FormFile("tweak")
		if err != nil {
			response := shared.BuildResponse(1, nil, fmt.Sprintf("File is not valid! Error: %s", err.Error()), w)
			w.Write(response)
			return
		}

		defer file.Close()
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			response := shared.BuildResponse(1, nil, fmt.Sprintf("File is not valid! Error: %s", err.Error()), w)
			w.Write(response)
			return
		}

		fileType := http.DetectContentType(fileData)
		if fileType == "application/zip" {

			zipName := shared.GenID()
			newPath := filepath.Join(shared.TempDir, zipName+".zip")
			newFile, err := os.Create(newPath)
			if err != nil {
				response := shared.BuildResponse(1, nil, fmt.Sprintf("Could not write to TEMP directory! Error: %s", err.Error()), w)
				w.Write(response)
				return
			}

			defer newFile.Close()
			if _, err := newFile.Write(fileData); err != nil || newFile.Close() != nil {
				response := shared.BuildResponse(1, nil, fmt.Sprintf("Could not write data to file! Error: %s", err.Error()), w)
				w.Write(response)
				return
			}

			if shared.IsMacOS() {
				macOS := new(bridge.MacOS)
				err := macOS.AddTweakPlugin(newPath)
				if err != nil {
					response := shared.BuildResponse(1, nil, fmt.Sprintf("Could not read zip from TEMP! Error: %s", err.Error()), w)
					w.Write(response)
					return
				}

				response := shared.BuildResponse(0, "Tweak has been installed!", nil, w)
				w.Write(response)
				return
			} else if shared.IsWindows() {
				/*windows := new(bridge.Windows)
				err := windows.AddTweakPlugin(newPath)
				if err != nil {
					response := shared.BuildResponse(1, nil, fmt.Sprintf("Could not read zip from TEMP! Error: %s", err.Error()), w)
					w.Write(response)
					return
				}*/
				//TODO: add windows implementation
				response := shared.BuildResponse(1, nil, "A Windows code injection runtime is in the works! Currently unsupported!", w)
				w.Write(response)
				return
			} else {
				response := shared.BuildResponse(1, nil, fmt.Sprintf("Operating System unsupported! Detected OS: %s", runtime.GOOS), w)
				w.Write(response)
				return
			}
		} else {
			response := shared.BuildResponse(1, nil, fmt.Sprintf("Invalid tweak! Must be ZIP format! Detected Type: %s", fileType), w)
			w.Write(response)
			return
		}

	case "remove":

		removalData := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&removalData)
		if err != nil {
			response := shared.BuildResponse(1, nil, fmt.Sprintf("Could not understand response from /remove! Error: %s", err.Error()), w)
			w.Write(response)
			return
		}

		if id, ok := removalData["pkg_id"]; ok {
			if id != "" {
				macOS := new(bridge.MacOS)
				err := macOS.RemoveTweakPlugin(id)
				if err != nil {
					response := shared.BuildResponse(1, nil, fmt.Sprintf("Could not remove tweak! Error: %s", err.Error()), w)
					w.Write(response)
					return
				}
			} else {
				response := shared.BuildResponse(1, nil, "Could not remove tweak! Identifier was left empty!", w)
				w.Write(response)
				return
			}
		} else {
			response := shared.BuildResponse(1, nil, "Package ID could not be found!", w)
			w.Write(response)
			return
		}

		response := shared.BuildResponse(0, "Finished removing tweak!", nil, w)
		w.Write(response)
		return
	default:
		response := shared.BuildResponse(1, nil, fmt.Sprintf("Action %s was not defined!", action), w)
		w.Write(response)
		return
	}
}

//NotFoundHandler - 404 route
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := shared.BuildResponse(1, nil, "Could not found route!", w)
	logger.Info(string(response))
	w.Write(response)
}
