package lang

type SoftwareLanguageInspector interface {
	IsInstalled() bool
	GetVersion() (SoftwareVersion, error)
	GetLatestVersion() (SoftwareVersion, error)
}

type SoftwareLanguageManager interface {
	GetPackageManagers()
	SetPackageManager(pm PackageManagerName, options ...OptionerFn) error
	SoftwareLanguageInspector
}

type PackageManManager interface {
	GetName() PackageManagerName
	GetLanguage() SoftwareLanguageName
}

type PythonInstallation struct {
	version   SoftwareVersion
	latest    SoftwareVersion
	installed bool
}

func NewSoftwareLanguage(name SoftwareLanguageName, options ...OptionerFn) SoftwareLanguageManager {
	var inst SoftwareLanguageManager
	switch name {
	case "python":
		{
			inst := PythonInstallation{}
		}
	}
}

func (p *PythonInstallation) IsInstalled() bool {

}

func (p *PythonInstallation) GetVersion() (SoftwareVersion, error) {
	if p.IsInstalled() {

	}
}
