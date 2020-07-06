package ini

import (
	"github.com/itning/DouBanReptile/internal/error2"
	"github.com/itning/DouBanReptile/internal/preference"
	"gopkg.in/ini.v1"
	"os"
)

var defaultConfigFileName = "DouBanConfig.ini"

type PreferenceConfig interface {
	Write(*preference.Preference)
	Read() *preference.Preference
}

type Config struct {
}

func (c Config) Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (c Config) CreateFile() {
	file, err := os.Create(defaultConfigFileName)
	handlerError(err)
	defer func() {
		if err := file.Close(); err != nil {
			handlerError(err)
		}
	}()
}
func (c Config) Write(preference *preference.Preference) {
	if !c.Exists(defaultConfigFileName) {
		c.CreateFile()
	}
	cfg := ini.Empty()
	e := ini.ReflectFrom(cfg, preference)
	handlerError(e)
	err := cfg.SaveTo(defaultConfigFileName)
	handlerError(err)
}

func (c Config) Read() *preference.Preference {
	if !c.Exists(defaultConfigFileName) {
		c.CreateFile()
	}
	cfg, err := ini.Load(defaultConfigFileName)
	handlerError(err)
	p := new(preference.Preference)
	e := cfg.MapTo(p)
	handlerError(e)
	return p
}

func handlerError(e error) bool {
	if nil == e {
		return false
	} else {
		error2.GetImpl().Handler(e)
		return true
	}
}
