package file

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
	io.Writer
}

type Files []File

func (fs Files) Close() error {
	var err error
	for i := range fs {
		if fs[i] == nil {
			continue
		}
		err1 := fs[i].Close()
		if err1 != nil {
			err = err1
		}
	}
	return err
}

var memBuf = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func FileStream(r io.Reader, memThreshold int64) (sf File, fsize int64, err error) {
	b := memBuf.Get().(*bytes.Buffer)
	b.Reset()
	var w io.Writer = b

	n, err := io.CopyN(w, r, memThreshold+1)
	if err != nil && err != io.EOF {
		memBuf.Put(b)
		return
	}
	fsize = n

	if n > memThreshold {
		f, err := ioutil.TempFile("", "")
		if err != nil {
			memBuf.Put(b)
			return nil, 0, err
		}
		sf = &FileCloser{f}

		_, err = io.Copy(f, b)
		memBuf.Put(b)
		if err != nil {
			_ = sf.Close()
			return nil, 0, err
		}

		n, err = io.Copy(f, r)
		if err != nil {
			_ = sf.Close()
			return nil, 0, err
		}
		fsize += n

		_, err = f.Seek(0, 0)
		if err != nil {
			_ = sf.Close()
		}
		return sf, fsize, err
	}

	sf = bytesFile{
		Reader: bytes.NewReader(b.Bytes()),
		closer: func() error {
			memBuf.Put(b)
			return nil
		},
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func newBytesReader(data []byte) File {
	b := memBuf.Get().(*bytes.Buffer)
	b.Reset()

	return bytesFile{
		Reader: bytes.NewReader(data),
		closer: func() error {
			memBuf.Put(b)
			return nil
		},
	}
}

type bytesFile struct {
	*bytes.Reader
	closer func() error
}

func (rc bytesFile) Close() error {
	return rc.closer()
}

func (rc bytesFile) Write(b []byte) (n int, err error) {
	err = errors.New("bytesFile not implemented")
	return
}

func newSectionFiles(fh File, sectionSize int64, sectionCount int) []File {
	var fhs []File
	for i := 0; i < sectionCount; i++ {
		fh := sectionFile{
			fh: fh,
			SectionReader: io.NewSectionReader(fh, int64(i)*sectionSize, sectionSize),
		}
		fhs = append(fhs, fh)
	}
	return fhs
}

type sectionFile struct {
	fh File
	*io.SectionReader
	closer func() error
}

func (sr sectionFile) Close() error {
	_ = sr.fh.Close()
	return nil
}

func (sr sectionFile) Write(b []byte) (n int, err error) {
	err = errors.New("sectionFile not implemented")
	return
}

func NewTmpFiles(n int) (fhs []File, err error) {
	for i := 0; i < n; i++ {
		fh, err := NewTmpFile()
		if err != nil {
			return nil, err
		}
		fhs = append(fhs, fh)
	}
	return
}

func NewTmpFile() (fh File, err error) {
	tmpFh, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}

	fh = &FileCloser{tmpFh}
	return
}

type FileCloser struct {
	*os.File
}

func (r *FileCloser) Close() error {
	name := r.File.Name()
	err := r.File.Close()
	if err != nil {
		return err
	}
	return os.Remove(name)
}
