package i18n

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	stlslices "github.com/kkkunny/stl/container/slices"
	"gopkg.in/yaml.v3"

	"github.com/kkkunny/pokemon/src/config"
	"github.com/kkkunny/pokemon/src/consts"
)

var translations map[string]string

func loadLocalisationFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	locs := make(map[string]string)
	err = decoder.Decode(&locs)
	if err != nil {
		return err
	}
	for k, v := range locs {
		translations[k] = v
	}
	return nil
}

func loadLocalisation(lang consts.Language) error {
	dirpath := filepath.Join(consts.LocalisationPath, string(lang))
	dirinfo, err := os.Stat(dirpath)
	if err != nil {
		return err
	} else if !dirinfo.IsDir() {
		return fmt.Errorf("%s is not a localisation directory", dirpath)
	}

	return filepath.WalkDir(dirpath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".yml" {
			return nil
		}
		err = loadLocalisationFile(path)
		if err != nil {
			return err
		}
		return nil
	})
}

func Init() error {
	translations = make(map[string]string)
	return loadLocalisation(config.Language)
}

func Get(key string, defaultValue ...string) string {
	s, ok := translations[key]
	if ok {
		return s
	}
	return stlslices.Last(defaultValue)
}
