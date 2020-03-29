package bridge

//Windows - windows type
type Windows int

//AddTweakPlugin - creates a new tweak by
func (os *Windows) AddTweakPlugin(fileName string) error {
	return nil
}

//RemoveTweakPlugin - removes a tweak by modifier id
func (os *Windows) RemoveTweakPlugin(pkgID string) error {
	return nil
}
