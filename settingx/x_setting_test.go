package settingx

import (
	"log"
	"os"
	"testing"
)

func TestSetting(t *testing.T) {
	Setup(os.Getwd())
	var appSetting struct {
		Env string
	}
	MapTo("app", &appSetting)
	log.Println(appSetting.Env)
}
