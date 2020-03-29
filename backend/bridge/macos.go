package bridge

import "webject/shared"

//consts for macos
const (
	InstallDirName = "/Library/macStr8t/Plugins"
)

//MacOS - macos type
type MacOS int

//AddTweakPlugin - creates a new tweak by
func (os *MacOS) AddTweakPlugin(fileName string) error {
	err := shared.Decompress(fileName, InstallDirName)
	if err != nil {
		return err
	}

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
