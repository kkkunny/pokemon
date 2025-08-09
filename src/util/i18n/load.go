package i18n

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tnnmigga/enum"
	"gopkg.in/yaml.v3"

	"github.com/kkkunny/pokemon/src/config"
)

type Language string

var LanguageEnum = enum.New[struct {
	ZH_CN Language `enum:"zh_cn"`
}]()

func loadLocalisationFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	locs := make(map[string]string)
	err = decoder.Decode(&locs)
	if err != nil {
		return nil, err
	}
	return locs, nil
}

func LoadLocalisation(lang Language) (*Localisation, error) {
	dirpath := filepath.Join(config.LocalisationPath, string(lang))
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	} else if !dirinfo.IsDir() {
		return nil, fmt.Errorf("%s is not a localisation directory", dirpath)
	}

	loc := NewLocalisation()
	err = filepath.WalkDir(dirpath, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) != ".yml" {
			return nil
		}
		kvs, err := loadLocalisationFile(path)
		if err != nil {
			return err
		}
		loc.MultiAdd(kvs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return loc, nil
}
