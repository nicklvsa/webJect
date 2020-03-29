package bridge

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"webject/shared"
)

//consts for macos
const (
	InstallDirName = "/Library/macSubstrate/Plugins"
)

var (
	logger *shared.Logger
)

//MacOS - macos type
type MacOS int

//AddTweakPlugin - creates a new tweak by
func (macos *MacOS) AddTweakPlugin(fileName string) error {
	err := shared.Decompress(fileName, InstallDirName)
	if err != nil {
		return err
	}

	err = shared.CleanDecompression(fileName)
	if err != nil {
		return err
	}

	return nil
}

//RemoveTweakPlugin - removes a tweak by modifier id
func (macos *MacOS) RemoveTweakPlugin(pkgID string) error {
	err := shared.RemoveWalkDirs(InstallDirName + string(os.PathSeparator) + pkgID)
	if err != nil {
		return err
	}
	return nil
}

//GetBundleIdentifierByApp - returns a mac app's bundle identifier by its path
func (macos *MacOS) GetBundleIdentifierByApp(appName string) (string, error) {

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd := exec.Command("osascript", "-e", fmt.Sprintf("'id of app \"%s\"'", appName))
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	fmt.Println(cmd.String())

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("STD_ERR: %s || GO_ERR: %s", stdErr.String(), err.Error())
	}

	return stdOut.String(), nil
}
