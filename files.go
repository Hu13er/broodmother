package broodmother

import (
	"io/ioutil"
	"os"
	"path"
)

type File struct {
	Path    string
	Content string
}

func (f *File) write() error {
	d, _ := path.Split(f.Path)
	err := os.MkdirAll(d, 0o700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.Path, []byte(f.Content), 0o700)
}
