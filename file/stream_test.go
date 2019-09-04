package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func Test_FileStream(t *testing.T) {
	fh, err := os.Open("test.txt")
	if err != nil {
		t.Error(err)
	}

	fi, err := fh.Stat()
	if err != nil {
		t.Error(err)
	}

	fCheck := func(fh *os.File, memSize int64) {
		f, fsize, err := FileStream(fh, memSize)
		if err != nil && err != io.EOF {
			t.Error(err)
		}

		if fsize != fi.Size() {
			t.Error(err)
		}

		_, err = fh.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		fileContent, err := ioutil.ReadAll(fh)
		if err != nil {
			t.Error(err)
		}

		streamContent, err := ioutil.ReadAll(f)
		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(fileContent, streamContent) != 0 {
			t.Error("stream content not match")
		}
	}

	fCheck(fh, fi.Size())
	_, err = fh.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	fCheck(fh, fi.Size()-1)
}
