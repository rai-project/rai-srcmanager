package srcmanager

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/repositories": {
		local:   "repositories",
		size:    2708,
		modtime: 1510767947,
		compressed: `
H4sIAAAAAAAA/4yWQY7kOgiG932Nt85rqTdzgbnHiNjEoQobC+NUVZ9+lG7N7P5odpE+Ygz5QCkS+1z/
T1bfnWTpbjdO8U6edjn4DeEZO2SPgdDqdmeHdNZ+sA+xhkISbRu81Bf8wNSrQqjCLSCtGaIWd8isbVIu
aPAT5zTni/v2CVFIxT1OM5O0zRDPFLTSgD3OekYMhvfOujlVfpjDvmRLFxZ8002Ul3WKZhz5828kFC7b
o6lRrtSo4JO4ksBe80E6KS6sLLZcivmHQzfPgNlDLrjKeusMbSq21GfjWLpzlhQGay22dOoXqU5tMc0W
3A7MrSgntQknpmB1i/eE2M6keOdIpQLbL00h05wo7RjLeofFqiXouFq58K1S2qXBrJXDJUGpq2WGtlZr
EubSoCpfniDYLPMvfnbzwNdvh2ShfwxbRoWydXJSxcV09s28UkuwV10660Uvu0rD+d3O9QHxXMdcEXWC
557P33tsOUxnhRnO5zu/CsPd4iRwlJyCVarAr+lMKZZktVvjFtAo5yIj/IW5Mg1ewkzhISM5RYIjOjg5
9m6wC6l8YpsG+3FNBVsywvwCOlOF9DWC69WsB7dhvqk9YIR52n+8/YewU+Le4dsnhpXPEPxN5sT2HNLZ
gwp8+ZBPhB4Mx+J5+4AntqPqUvDPhyq92N9+BwAA//8a+tFBlAoAAA==
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
