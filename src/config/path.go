package config

import (
	"path/filepath"

	stlerr "github.com/kkkunny/stl/error"
	stlos "github.com/kkkunny/stl/os"
)

var (
	RootPath         = string(stlerr.MustWith(stlos.GetWorkDirectory()))
	ResourcePath     = filepath.Join(RootPath, "resource")
	FontsPath        = filepath.Join(ResourcePath, "fonts")
	LocalisationPath = filepath.Join(ResourcePath, "localisation")
	MapItemPath      = filepath.Join(ResourcePath, "map_item")
	VoicePath        = filepath.Join(ResourcePath, "voice")
	ScriptsPath      = filepath.Join(ResourcePath, "scripts")
)
