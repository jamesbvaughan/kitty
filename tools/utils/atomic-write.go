// License: GPLv3 Copyright: 2022, Kovid Goyal, <kovid at kovidgoyal.net>

package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var _ = fmt.Print

func AtomicWriteFile(path string, data []byte, perm os.FileMode) (err error) {
	npath, err := filepath.EvalSymlinks(path)
	if errors.Is(err, fs.ErrNotExist) {
		err = nil
		npath = path
	}
	if err == nil {
		path = npath
		path, err = filepath.Abs(path)
		if err == nil {
			var f *os.File
			f, err = os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".atomic-write-")
			if err == nil {
				removed := false
				defer func() {
					f.Close()
					if !removed {
						os.Remove(f.Name())
						removed = true
					}
				}()
				_, err = f.Write(data)
				if err == nil {
					err = f.Chmod(perm)
					if err == nil {
						err = os.Rename(f.Name(), path)
						if err == nil {
							removed = true
						}
					}
				}
			}
		}
	}
	return
}

func AtomicUpdateFile(path string, data []byte, perms ...fs.FileMode) (err error) {
	perm := fs.FileMode(0o666)
	if len(perms) > 0 {
		perm = perms[0]
	}
	s, err := os.Stat(path)
	if err == nil {
		perm = s.Mode().Perm()
	}
	return AtomicWriteFile(path, data, perm)
}
