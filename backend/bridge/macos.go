package bridge

import (
	"fmt"
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
func (os *MacOS) AddTweakPlugin(fileName string) error {
	files, err := shared.Decompress(fileName, InstallDirName)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Files Unzipped: %s", files))

	err = shared.CleanDecompression(fileName, InstallDirName)
	if err != nil {
		return err
	}

	return nil
}

//RemoveTweakPlugin - removes a tweak by modifier id
func (os *MacOS) RemoveTweakPlugin(pkgID string) error {
	return nil
}
