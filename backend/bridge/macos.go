package bridge

import (
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
func (macos *MacOS) GetBundleIdentifierByApp(appPath string) (string, error) {
	out, err := exec.Command("mdls", "-name kMDItemCFBundleIdentifier", fmt.Sprintf("-r %s", appPath)).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
