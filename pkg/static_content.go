// Code generated by "esc -o pkg/static_content.go -pkg srcmanager -private repositories"; DO NOT EDIT.

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

	"/repositories": {
		local:   "repositories",
		size:    3638,
		modtime: 1540993451,
		compressed: `
H4sIAAAAAAAC/4yXQc7jNgyF9znFALN2B/g3vUBvURQBLdEOE0oUKNpJ5vSF86PdPSXr74WSqccn5ecq
cdnmP5KVX04yNbcrp/gVTolbOwFMni6yM8RbXCC7d4Rmtxs7pFtpO3sXq0iSaFl4CL8QLVoss/ZkDVdQ
4RqQlgxRjRtkVhdZBzT4gdc056KQtg2ikIJbnbZMUhdDPFPQTB02Kuuh6BxYsDgVvpvfsETpibeYLd3e
0UWUp3kTzVj51//KjovdqxrlQpVWXIkLCTwK3kk3ioF3V5uG9v2Pf40EA58deGshA64yXxuvA0V5VI6p
OWdJYT5QNspZ+cffzWWn4H+G0jbaVXDt5h5YcowQptmC6465rcpJbYPTu+IxWr0lxC5MimNQCq3wrKUq
ZJoTpQvGMt/gx6qlG2brwNyF0kUqXLVwuCQ4QUWS28w1wXa8whfDKmEuFVrz5UsEq2U+86OZB/7AuksW
+lA29SJYUXRaYXR+zwWmTqq4E419MS9UEy4hjXVwVE2lws03tyMKId7mvs2IOskATd+ZPO2mW+GR8MbP
letAkTELVikSp58/kIQpxZSsNKtcA1rWeZUe/hxVUqbOU5gpLNOTU2Dbd06OjdvZhVR+Yzt29n1MBTul
h/kAOlOB9NmDyyhOvlN7UbuPFTjXwzxd/nx/g7yeqrAHWwg+nW3DXtqlsQetHQt+v9/cneG4PK5fODB1
Ik/T4NfBqxPOoGOQEV3aSoP9tsHFdDwMz9TkPG56wvl35UKqlgYr1Pqm/PFmeL7sgRRW6wNnbO8qM74H
Kc+kw4g9/oX04BofuPPl8uOdfn5dgJ/+wlp/L816XAjnO89dgj/Vv9f1J3QmJT2d/g0AAP//+N9pAjYO
AAA=
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
