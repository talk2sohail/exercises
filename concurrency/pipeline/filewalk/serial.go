package filewalk

import (
	"crypto/md5"
	"os"
	"path/filepath"
)

func Serail_MD5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}
