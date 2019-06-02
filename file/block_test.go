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
	fname := "/Users/wf/Downloads/1.mp4"

	//test file EC shards
	//shards := blockMgr.DataShards + blockMgr.ParShards
	//rds := make([]io.Reader, shards)

	//ok, err := blockMgr.Verify(rds)
	//assert.Equal(t, true, ok)
	//assert.Equal(t, nil, errhao)

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

func TestSplit(t *testing.T) {
	fname := "/Users/wf/Downloads/test1.txt"
	f, err := os.Open(fname)
	assert.Equal(t, nil, err)
	defer f.Close()

	blockMgr, err = NewBlockMgr(64, 32)
	if err != nil {
		panic(err)
	}

	fi, err := f.Stat()
	start := time.Now()

	fhs := make([]*os.File, blockMgr.DataShards)
	doneInex := make(chan int, blockMgr.DataShards)
	go blockMgr.Split(f, fhs, fi.Size(), doneInex)
	assert.Equal(t, nil, err)
	for {
		select {
		case i, ok := <-doneInex:
			if ok {
				fhs[i].Close()
				fmt.Println("====>", i)
			} else {
				fmt.Println("====> done time cost:", time.Now().Sub(start))
				return
			}
		}
	}
	return
}

func TestCalCid(t *testing.T) {
	fname := "/Users/wf/Downloads/test.txt"
	f, err := os.Open(fname)
	assert.Equal(t, nil, err)
	defer f.Close()

	n := time.Now()
	cid, err := GetCID(f)
	fmt.Println("========> cal cid cost:", time.Now().Sub(n))
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", cid)
}
