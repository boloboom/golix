package settingx

import (
	"log"
	"path"

	"gopkg.in/ini.v1"
)

var cfg *ini.File

// Setup initializes the application configuration.
//
// This function loads the 'app.ini' file and parses its contents.
// It sets the configuration object 'cfg' to the loaded configuration.
// If there is an error while loading or parsing the file, it logs a fatal error.
//
// No parameters are taken.
// No return values.
func Setup(dir string, err error) {
	if err != nil {
		log.Fatalf("setting.Setup, os.Getwd() failed: %v", err)
	}
	settingPath := path.Join(dir, "app.ini")
	cfg, err = ini.Load(settingPath)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse %s: %v", settingPath, err)
	}
}

// MapTo maps the configuration section to the provided interface.
//
// It takes a string parameter 'section' which represents the name of the configuration section to be mapped.
// It also takes an interface{} parameter 'v' which represents the target struct or variable to which the configuration section will be mapped.
// The function does not return anything.
func MapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
