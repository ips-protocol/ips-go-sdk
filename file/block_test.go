package file

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var blockMgr *BlockMgr

func init() {
	var err error
	blockMgr, err = NewBlockMgr(64, 32)
	if err != nil {
		panic(err)
	}
	return
}

func TestBlockMgr_SplitFile(t *testing.T) {
	fname := "test.txt"

	//test stean EC shards
	f, err := os.Open(fname)
	assert.Equal(t, nil, err)
	defer f.Close()
	fi, err := f.Stat()
	assert.Equal(t, nil, err)

	rs, err := blockMgr.ECShards(f, fi.Size())
	assert.Equal(t, err, nil)

	ok, err := blockMgr.Verify(rs)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	//test shards join
	f.Seek(0, io.SeekStart)

	rs2, _ := blockMgr.ECShards(f, fi.Size())
	fileContent, err := ioutil.ReadAll(f)
	assert.Equal(t, nil, err)
	fileContent2 := bytes.NewBuffer(nil)
	blockMgr.Join(fileContent2, rs2, fi.Size())
	bytes.Equal(fileContent2.Bytes(), fileContent)
}

func TestCalCid(t *testing.T) {
	fname := "test.txt"
	f, err := os.Open(fname)
	assert.Equal(t, nil, err)
	defer f.Close()

	n := time.Now()
	cid, err := GetCID(f)
	fmt.Println("========> cal cid cost:", time.Now().Sub(n))
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", cid)
}

func TestSplit(t *testing.T) {
	fp := os.Getenv("TEST_SPLIT_FILE_PATH")
	if fp == "" {
		return
	}
	fh, err := os.Open(fp)
	assert.NoError(t, err)

	fi, err := os.Stat(fp)
	assert.NoError(t, err)

	dataShards, parShards, _ := BlockCount(fi.Size())
	mgr, err := NewBlockMgr(dataShards, parShards)
	assert.NoError(t, err)

	fhs1, err := mgr.Split(fh, fi.Size())
	assert.NoError(t, err)

	fhs2, err := mgr.Split2(fh, fi.Size())
	assert.NoError(t, err)
	assert.Equal(t, len(fhs1), len(fhs2))

	defer func() {
		Files(fhs1).Close()
		Files(fhs2).Close()
	}()

	for i := range fhs1 {
		h1 := md5.New()
		_, err := io.Copy(h1, fhs1[i])
		assert.NoError(t, err)
		h1Md5 := h1.Sum(nil)

		h2 := md5.New()
		_, err = io.Copy(h2, fhs2[i])
		assert.NoError(t, err)
		h2Md5 := h2.Sum(nil)

		assert.Equal(t, h1Md5, h2Md5)
	}

}
