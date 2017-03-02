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

	"/LICENSE.TXT": {
		local:   "../LICENSE.TXT",
		size:    1735,
		modtime: 1488460929,
		compressed: `
H4sIAAAAAAAA/7RUy27juBJdX37FQVbdgK7nsZhF7xSLjomxKUOiO5MlLZUjTiRRIKkE/vsBqcTtnqwn
mwClOnUeVfTaThdnnruAL81X/P7rb39kyE/t3KPQLy/6BXnfIzV4OPLkXqldMVbQK/V2ohanyzf87xbB
8K8/sT/ka4UHZ+fp08fjaF7JeRMusGeIvjejNT7D0Z30qP+/7vQwafM8fgKaYdJNWDWua1bmHbaidmbs
QG4w3hs7AsYDHTk6XYBnp8dAbQacHREiH5pOu2fKgGABPV7YRM5HoD0FbUYzPkMDjZ0uqT90caC35/Cm
HUGPLaC9t43RgVrW2mYeaAw6JPKz6ckDX0JHwF39jrr7miW2lnQPvJnQAaEj9vE9lewcEPMOzjRxWAYz
Nv3cRkHXht4MZuFikeJ9S8Fi9pQl0RkG25pz/E/J5zSfeuO7DK2Js09zoIz5WGxojKjo6RfrAE99H2eY
aCFGlVx8iMwW78HGwAYTsOSW2PHW2SG1Xy0Zj/PsRuM7Sii0Ft4m5r+pCbES+8+27+1b9NjYsTXRmf/G
WEVXtbGyyPF2dg3F3bQEYJh9Sixu7V2rPtlXYs31wEcbTEPZssTexH57vmFKln6SwVrjm16bgZxffZZh
RpzMqN0FOFs3/BAxOdvODf0XMhLrMvjna/vYmw0dOWDQgZzRvWeTs6+mpXY5tYi8tbFikkyCpJmjHsin
fG+fdAaI9NyWR5wBo00Adu03wUcPy1zrPAZ9AU4Ub7GN+6Wxtc4TrMPk7GADsSWm4NGSM6/xMs4uXY7x
+Pwa/ESNOZsGkzPWsTdnQqAR0/W1r5jacqAuN+oxrzggauBQld9FwQvgLq8h6rsMeBRqWx4V8JhXVS7V
E1BuWC6fgD+FLDLwvw4Vr2ugrNLP107wIoOQ692xEPIBuD8qyFIBO7EXihdMlUBi/5gpeIRvgD2v1ttc
qvxe7IR6yoCNUDJNxybOz9khr5RYH3d5hcOxOpQ1B3JZQJZSyE0l5APfc6lWgJCQJcC/c6lQb/PdbvG8
LqWqxP1RlVWdRK/Lw1MlHrYK2Ja7glc17nmUm9/v+DuvfGLrXS72GYp8nz/wBVmqLa9S47vexy1PJSEj
BvlaiVLGwBJpvlYZoMpK/UA/ippnyCtRp7A2VbnPkAIvN0ukMoIlX0bFdSzpXVcXm2LhWPMbUQXPd0I+
1CwhbgGrfwIAAP//d3yIgscGAAA=
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
