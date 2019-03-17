package file

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

var blockMgr *BlockMgr

func init() {
	cfg := Config{4, 2}
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

	rs, err := blockMgr.ECShards(f, fi.Size())
	ok, err = blockMgr.Verify(rs)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

}
