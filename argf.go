// Package argf implements an io.Reader similar to ruby's ARGF
package argf

import (
	"fmt"
	"io"
	"os"
)

func Name(r io.Reader) string {
	return ioReaderName(r)
}

// Argf returns an io.Reader from STDIN or files from command-line arguments.
// If there are multiple files, they are concatinated as one io.Reader. A
// filename of "-" refers to stdin.
func Argf() io.Reader {
	r, err := FromNames(os.Args[1:]...)
	if err != nil {
		panic(err)
	}
	return r
}

func ioReaderName(r io.Reader) string {
	if r == os.Stdin {
		return ":stdin:"
	}

	if mr, ok := r.(*multiReader); ok {
		if len(mr.readers) > 0 {
			return ioReaderName(mr.readers[0])
		}
		return ":EOF:"
	}

	if f, ok := r.(*os.File); ok {
		return f.Name()
	}

	return fmt.Sprintf("%T", r)
}

// FromNames returns an io.Reader from STDIN or files from function argument.
func FromNames(filenames ...string) (io.Reader, error) {
	// https://raw.githubusercontent.com/yuya-takeyama/argf/master/argf.go
	var reader io.Reader

	filelen := len(filenames)

	if filelen > 0 {
		files := make([]io.Reader, filelen)

		for i := 0; i < filelen; i++ {
			if filenames[i] == "-" {
				files[i] = os.Stdin
				continue
			}

			file, err := os.Open(filenames[i])
			if err != nil {
				return nil, err
			}

			files[i] = file
		}

		r := make([]io.Reader, len(files))
		copy(r, files)
		reader = &multiReader{r}
	} else {
		reader = os.Stdin
	}

	return reader, nil
}

func From(readers ...io.Reader) (io.Reader, error) {
	return &multiReader{readers}, nil
}

// from io.MultiReader

type multiReader struct {
	readers []io.Reader
}

func (mr *multiReader) Read(p []byte) (n int, err error) {
	for len(mr.readers) > 0 {
		// Optimization to flatten nested multiReaders (Issue 13558)
		if len(mr.readers) == 1 {
			if r, ok := mr.readers[0].(*multiReader); ok {
				mr.readers = r.readers
				continue
			}
		}
		n, err = mr.readers[0].Read(p)
		if err == io.EOF {
			mr.readers = mr.readers[1:]
		}
		if n > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				// Don't return io.EOF yet. More readers remain.
				err = nil
			}
			return
		}
	}
	return 0, io.EOF
}
