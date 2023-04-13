package fs

import (
	"io/fs"
	"path/filepath"
)

func FindFilesWithExtension(root, ext string) []string {
	var files []string
	filepath.WalkDir(root, func(file string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(dir.Name()) == ext {
			files = append(files, file)
		}
		return nil
	})
	return files
}
