//Package filesystem provides filesystem based Stages for Slurp.

package fs

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/omeid/slurp"
	"github.com/omeid/slurp/tools/glob"
)

// A simple helper function that reads the file from the given path and
// returns a pointer to a slurp.File or an error.
func Read(path string) (*slurp.File, error) {
	Stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fs := &slurp.File{Reader: f, Path: path, FileInfo: slurp.FileInfoFrom(Stat)}

	return fs, nil
}

//Src returns a channel of slurp.Files that match the provided pattern.
func Src(c *slurp.C, globs ...string) slurp.Pipe {

	pipe := make(chan slurp.File)

	files, err := glob.Glob(globs...)

	if err != nil {
		c.Error(err)
		close(pipe)
	}

	cwd, err := os.Getwd()
	if err != nil {
		c.Error(err)
		close(pipe)
		return pipe
	}

	//TODO: Parse globs here, check for invalid globs, split them into "filters".
	go func() {
		defer close(pipe)

		for matchpair := range files {

			f, err := Read(matchpair.Name)
			if err != nil {
				c.Error(err)
				continue
			}

			f.Cwd = cwd
			f.Dir = glob.Dir(matchpair.Glob)
			pipe <- *f
		}

	}()

	return pipe
}

// Dest writes the files from the input channel to the dst folder and closes the files.
// It never returns Files.
func Dest(c *slurp.C, dst string) slurp.Stage {
	return func(files <-chan slurp.File, out chan<- slurp.File) {

		var wg sync.WaitGroup
		defer wg.Wait()

		for file := range files {

			realpath, _ := filepath.Rel(file.Dir, file.Path)
			path := filepath.Join(dst, filepath.Dir(realpath))
			err := os.MkdirAll(path, 0700)
			if err != nil {
				c.Error(err)
				return
			}

			if !file.FileInfo.IsDir() {

				wg.Add(1)
				go func(file slurp.File) {

					defer wg.Done()
					defer file.Close()

					realfile, err := os.Create(filepath.Join(dst, realpath))
					if err != nil {
						c.Error(err)
						return
					}
					io.Copy(realfile, file)
					realfile.Close()

				}(file)
			}
		}

	}
}
