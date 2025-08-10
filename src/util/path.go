package util

import (
	"os"
	"path/filepath"

	stlval "github.com/kkkunny/stl/value"
)

// 在目录下查找指定名字（不包含拓展名）的文件并执行parse函数
func FindFileAndThenParse[T any](dirpath, filename string, parseFn func(string) (T, error)) (T, error) {
	fs, err := os.ReadDir(dirpath)
	if err != nil {
		return stlval.Default[T](), err
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		curFileExt := filepath.Ext(f.Name())
		curFilenameWithoutExt := f.Name()
		if len(curFileExt) > 0 {
			curFilenameWithoutExt = curFilenameWithoutExt[:len(curFilenameWithoutExt)-len(curFileExt)]
		}
		if curFilenameWithoutExt != filename {
			continue
		}
		curFilepath := filepath.Join(dirpath, f.Name())
		return parseFn(curFilepath)
	}
	return stlval.Default[T](), os.ErrNotExist
}
