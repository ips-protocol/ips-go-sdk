package file

import (
	"bytes"
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
