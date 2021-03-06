// Code generated by "esc -o cmd/static_content.go -pkg cmd -private LICENSE.TXT"; DO NOT EDIT.

package cmd

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
		_ = f.Close()
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

	"/LICENSE.TXT": {
		local:   "LICENSE.TXT",
		size:    1736,
		modtime: 1519313858,
		compressed: `
H4sIAAAAAAAC/7RUy27juBJdX37FQVbdgK7nsRmgd4pFx8TYlCHRncmSlsoRJ5IokFQC//2AVOJ2T9aT
TYBSnTqPKnptp4szz13Al+Yrfv/1tz8y5Kd27lHolxf9grzvkRo8HHlyr9SuGCvolXo7UYvT5Rv+d4tg
+Nef2B/ytcKDs/P06eNxNK/kvAkX2DNE35vRGp/h6E561P9fd3qYtHkePwHNMOkmrBrXNSvzDltROzN2
IDcY740dAeOBjhydLsCz02OgNgPOjgiRD02n3TNlQLCAHi9sIucj0J6CNqMZn6GBxk6X1B+6ONDbc3jT
jqDHFtDe28boQC1rbTMPNAYdEvnZ9OSBL6Ej4K5+R919zRJbS7oH3kzogNAR+/ieSnYOiHkHZ5o4LIMZ
m35uo6BrQ28Gs3CxSPG+pWAxe8qS6AyDbc05/qfkc5pPvfFdhtbE2ac5UMZ8LDY0RlT09It1gKe+jzNM
tBCjSi4+RGaL92BjYIMJWHJL7Hjr7JDar5aMx3l2o/EdJRRaC28T89/UhFiJ/Wfb9/Ytemzs2JrozH9j
rKKr2lhZ5Hg7u4bibloCMMw+JRa39q5Vn+wrseZ64KMNpqFsWWJvYr893zAlSz/JYK3xTa/NQM6vPssw
I05m1O4CnK0bfoiYnG3nhv4LGYl1GfzztX3szYaOHDDoQM7o3rPJ2VfTUrucWkTe2lgxSSZB0sxRD+RT
vrdPOgNEem7LI86A0SYAu/ab4KOHZa51HoO+ACeKt9jG/dLYWucJ1mFydrCB2BJT8GjJmdd4GWeXLsd4
fH4NfqLGnE2DyRnr2JszIdCI6fraV0xtOVCXG/WYVxwQNXCoyu+i4AVwl9cQ9V0GPAq1LY8KeMyrKpfq
CSg3LJdPwJ9CFhn4X4eK1zVQVunnayd4kUHI9e5YCPkA3B8VZKmAndgLxQumSiCxf8wUPMI3wJ5X620u
VX4vdkI9ZcBGKJmmYxPn5+yQV0qsj7u8wuFYHcqaA7ksIEsp5KYS8oHvuVQrQEjIEuDfuVSot/lut3he
l1JV4v6oyqpOotfl4akSD1sFbMtdwasa9zzKze93/J1XPrH1Lhf7DEW+zx/4gizVllep8V3v45ankpAR
g3ytRCljYIk0X6sMUGWlfqAfRc0z5JWoU1ibqtxnSIGXmyVSGcGSL6PiOpb0rquLTbFwrPmNqILnOyEf
apYQt4AV+ycAAP//s0SJlsgGAAA=
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
