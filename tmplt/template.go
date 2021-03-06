package tmplt

import (
	htmltemplate "html/template"
	"io"
	"text/template"
)

// Common - in all templates
type Common struct {
	Title   string
	BaseURI string
}

//URLList - struct for url list type
type URLList struct {
	Data   [][]string
	Common Common
}

//List - struct for list type
type List struct {
	Data   []string
	Common Common
}

//TestData - struct for test type
type TestData struct {
	Path   string
	Host   string
	Common Common
}

//BootData - used for html boot config instructions
type BootData struct {
	User           string
	AppPath        string
	LogFile        string
	ConfigFile     string
	ConfigCheckURL string
	TestURL        string
	ServiceFile    string
	Common         Common
}

//ConfigData - used for default config file
type ConfigData struct {
	Port             string
	Path             string
	Cachepath        string
	Execplayer       string
	Awsconfigprofile string
	ForceBufferSize  string
}

//SysD - used for systemd file built for user
type SysD struct {
	User    string
	AppPath string
	AppName string
}

//UpnpDevDesc - UPNP Device Description data
type UpnpDevDesc struct {
	DeviceType   string
	FriendlyName string
	Manufacturer string
	ModelNumber  string
	SerialNumber string
	UUID         string
	Services     []UpnpService
}

//UpnpService - UPNP service structure
type UpnpService struct {
	Path string
}

//ParseHTM -- parses a template against base html template
func ParseHTM(wr io.Writer, tmpl string, data interface{}) error {
	b, err := htmltemplate.New("").Parse(baseHTML)
	if err != nil {
		return err
	}
	t, err := b.Parse(tmpl)
	if err != nil {
		return err
	}
	if err = t.Execute(wr, data); err != nil {
		return err
	}
	return nil
}

//ParseF -- parses a file template
func ParseF(wr io.Writer, tmpl string, data interface{}) error {
	b, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	if err = b.Execute(wr, data); err != nil {
		return err
	}
	return nil
}
