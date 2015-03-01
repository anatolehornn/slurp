package slurp

import (
	"io"
	"os"
	"time"
)


// File is the object that is passed between Slurp Pipes and Stages.
type File struct {
	io.Reader
	Cwd  string //Where are we?
	Dir  string //Dir, usually glob.Base
	Path string //Full path.

	FileInfo FileInfo
}

// Returns a copy of File.FileInfo.
// This method is provided for interoperability with os.File consider using 
// the following interface for a Slurp and os.File friendly API
//
//  type File interface {
//    io.ReadCloser
//    Stat() (os.FileInfo, error)
//  }
func (f File) Stat() (os.FileInfo, error) {
  return f.FileInfo, nil
}

func (f *File) Close() error {
	return Close(f.Reader)
}

func Close(in interface{}) error {
	if closer, ok := in.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func FileInfoFrom(fi os.FileInfo) FileInfo {
	return FileInfo{
		fi.Name(),
		fi.Size(),
		fi.Mode(),
		fi.ModTime(),
		fi.IsDir(),
		fi.Sys(),
	}

}

// Fileinfo implements os.FileInfo.
type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (f FileInfo) Name() string {
	return f.name
}
func (f FileInfo) Size() int64 {
	return f.size
}

func (f FileInfo) Mode() os.FileMode {
	return f.mode
}

func (f FileInfo) ModTime() time.Time {
	return f.modTime
}

func (f FileInfo) IsDir() bool {
	return f.isDir
}

func (f FileInfo) Sys() interface{} {
	return f.sys
}

func (f *FileInfo) SetName(name string) {
	f.name = name
}

func (f *FileInfo) SetSize(size int64) {
	f.size = size
}

func (f *FileInfo) SetMode(mod os.FileMode) {
	f.mode = mod
}

func (f *FileInfo) SetModTime(time time.Time) {
	f.modTime = time
}
func (f *FileInfo) SetIsDir(isdir bool) {
	f.isDir = isdir
}
func (f *FileInfo) SetSys(sys interface{}) {
	f.sys = sys
}
