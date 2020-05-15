package broodmother

type File struct {
	Path    string
	Content string
}

func (f *File) write() error {
	return nil
	// err := os.MkdirAll(path.Dir(f.Path), 0o700)
	// if err != nil {
	// 	return err
	// }
	// return ioutil.WriteFile(f.Path, []byte(f.Content), 0o700)
}
