package gin

import (
	"net/http"
	"os"
)

type (
	OnlyFilesFS struct {
		fs http.FileSystem
	}
	//Neutered  绝育 丧失某种能力
	NeuteredReaddirFile struct {
		http.File
	}
)

// It returns a http.Filesystem that can be used by http.FileServer(). It is used interally
// in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if listDirectory {
		return fs
	} else {
		return &OnlyFilesFS{fs}
	}

}

// Conforms to http.Filesystem
func (fs OnlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	return NeuteredReaddirFile{f}, nil
}

// Overrides the http.File default implementation
func (f NeuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
