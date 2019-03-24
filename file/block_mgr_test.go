package file

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var blockMgr *BlockMgr

func init() {
	cfg := Config{DataShards: 4, ParShards: 2}
	var err error
	blockMgr, err = NewBlockMgr(cfg)
	if err != nil {
		panic(err)
	}
	return
}

func TestBlockMgr_SplitFile(t *testing.T) {
	fname := "./test.txt"

	//test file EC shards
	rcs, err := blockMgr.SplitFile(fname)
	assert.Equal(t, nil, err)
	shards := blockMgr.DataShards + blockMgr.ParShards
	rds := make([]io.Reader, shards)
	for i := range rcs {
		rds[i] = rcs[i]
		defer rcs[i].Close()
	}

	ok, err := blockMgr.Verify(rds)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	//test stean EC shards
	f, err := os.Open(fname)
	assert.Equal(t, nil, err)
	defer f.Close()
	fi, err := f.Stat()
	assert.Equal(t, nil, err)

	meta := []byte("abcd")
	rs, err := blockMgr.ECShards(f, meta, fi.Size())
	ok, err = blockMgr.Verify(rs)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

	//test shards join
	f.Seek(0, io.SeekStart)

	rs2, _ := blockMgr.ECShards(f, nil, fi.Size())
	fileContent, err := ioutil.ReadAll(f)
	assert.Equal(t, nil, err)
	fileContent2 := bytes.NewBuffer(nil)
	blockMgr.Join(fileContent2, rs2, fi.Size())
	bytes.Equal(fileContent2.Bytes(), fileContent)
}
