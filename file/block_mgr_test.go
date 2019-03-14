package file

import (
	"github.com/stretchr/testify/assert"
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

	rs, err := blockMgr.SplitFile(fname)
	assert.Equal(t, nil, err)

	shards := blockMgr.DataShards + blockMgr.ParShards
	rds := make([]io.Reader, shards)
	for i := range rs {
		rds[i] = rs[i]
		defer rs[i].Close()
	}

	//rs, err := blockMgr.SplitFile2(fname)
	//assert.Equal(t, nil, err)

	ok, err := blockMgr.Verify(rs)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)

}
