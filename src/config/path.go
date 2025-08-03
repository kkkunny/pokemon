package config

import (
	"path/filepath"

	stlerr "github.com/kkkunny/stl/error"
	stlos "github.com/kkkunny/stl/os"
)

var (
	RootPath         = string(stlerr.MustWith(stlos.GetWorkDirectory()))
	DataPath         = filepath.Join(RootPath, "data")
	FontsPath        = filepath.Join(DataPath, "fonts")
	LocalisationPath = filepath.Join(DataPath, "localisation")
	WorldPath        = filepath.Join(DataPath, "world")
	MapsPath         = filepath.Join(WorldPath, "maps")
	MapItemPath      = filepath.Join(DataPath, "map_item")
	VoicePath        = filepath.Join(DataPath, "voice")
	ScriptsPath      = filepath.Join(DataPath, "scripts")
)
