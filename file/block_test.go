package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var blockMgr *BlockMgr

func init() {
	var err error
	blockMgr, err = NewBlockMgr(2, 1)
	if err != nil {
		panic(err)
	}
	return
}

func TestBlockMgr_SplitFile(t *testing.T) {
	fname := "/Users/wf/Downloads/1.mp4"

	//test file EC shards
	//shards := blockMgr.DataShards + blockMgr.ParShards
	//rds := make([]io.Reader, shards)

	//ok, err := blockMgr.Verify(rds)
	//assert.Equal(t, true, ok)
	//assert.Equal(t, nil, err)

	//test stean EC shards
	f, err := os.Open(fname)
	assert.Equal(t, nil, err)
	defer f.Close()
	fi, err := f.Stat()
	assert.Equal(t, nil, err)

	rs, err := blockMgr.ECShards(f, fi.Size())
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
