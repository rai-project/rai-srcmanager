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
		size:    3045,
		modtime: 1518194525,
		compressed: `
H4sIAAAAAAAA/4yWQW4bOwyG9zlFgKznBcjmXaC3KIqCI3FkxpQoUJyxndMXjtGim1/pboCPQ3Goj7Rf
isRpX/9LVl+dZOlu75ziNZwS9/4EMHk6ycEQ73GC7DIQWt3O7JDutR/sQ6yhkETbBov6hG+YelUIVbgF
pDVD1OIMmbVNyoQGX/GZ5jypt+8QhVTc47RnkrYZ4pmCVhqwx1nvEYNh3Vk3p8oXc9iXrEo3XGK2NJHk
QTdRXtZdNOPIb38ioY/ZLk2NcqVGBWfiSgKvgg/SnWIibbFl6u1vDtW9B0w8u+O9h0y4yvreGbpYbKnX
xrF05ywpDLai2NIpZ+Xn793loOAf09A+qyq4DXOHLhW7jxCm2YLbgbkV5aS2w+kteIyK94TYiUnx/pNK
Bd61NIVMc6J0wljWM/xYtQT9UCsTuSulkzR4auVwSXCCqiS3lVuC7aiWGQ5PtSZhLg2q+eklgs0y/+Rr
Nw/8ge2QLPSPYcuo0Nh2VF0KXJ2PucDUSRV3orNv5pVawimks06uqqs0WHx3u69CiPd17CuiTjDv/fmx
k5fDdK/whPvzmW+F4Z50EjipTsEqVeLp5RmFMKVYktVujVtAZZ2LjPDbLJMyDV7CTGGakZwCaz84ORZ3
sAupfGAdB/sxp4JNGWE+gc5UIb2N4DpbJ4+tvald5hF4r4d5Ov3/9S/I539U2IM9BN/OvmOXDunsQQW+
fMjH18VdGI7L9f0NL0xdyNMyeTu4OE3WxJAR3OKvCp9+BQAA//+SjS4e5QsAAA==
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
