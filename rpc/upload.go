package rpc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"

	"github.com/ipfs/go-ipfs/metafile"
	"github.com/ipweb-group/go-sdk/file"
)

func (c *Client) UploadWithPath(fpath string) (cid string, err error) {
	fname := path.Base(fpath)
	fh, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer fh.Close()

	fi, err := fh.Stat()
	if err != nil {
		return
	}

	cid, err = c.Upload(fh, fname, fi.Size())
	return
}

func (c *Client) Upload(rdr io.Reader, fname string, fsize int64) (cid string, err error) {
	dataShards, parShards, shardSize := file.BlockCount(fsize)
	mgr, err := file.NewBlockMgr(dataShards, parShards)
	if err != nil {
		return
	}

	h := sha256.New()
	_, err = h.Write([]byte(c.WalletPubKey))
	if err != nil {
		return
	}
	r := io.TeeReader(rdr, h)

	//dataFhs := make([]*os.File, dataShards)
	dataFhs, err := mgr.Split(r, fsize)
	if err != nil {
		return
	}
	cid, err = file.GetCidV0(h)
	if err != nil {
		return
	}

	parFhs, err := file.CreateTmpFiles(parShards)
	if err != nil {
		return
	}

	fmt.Println("===> shards:", dataShards, "parshards:", parShards, "shardsize:", shardSize)
	dataFhRdrs := make([]io.Reader, dataShards)
	for i := range dataFhRdrs {
		dataFhRdrs[i] = dataFhs[i]
	}
	parFhWtrs := make([]io.Writer, parShards)
	for i := range parFhWtrs {
		parFhWtrs[i] = parFhs[i]
	}
	err = mgr.Encode(dataFhRdrs, parFhWtrs)
	if err != nil {
		return
	}
	fhs := append(dataFhs, parFhs...)
	defer func() {
		err := file.DeleteTempFiles(fhs)
		if err != nil {
			log.Println("delete files failed:", err)
		}
	}()

	meta := metafile.NewMeta(fname, cid, fsize, uint32(dataShards), uint32(parShards))
	meta.WalletPubKey = c.WalletPubKey
	shardSize += int64(len(meta.Encode(0)))
	shards := dataShards + parShards
	_, err = c.NewUploadJob(cid, fsize, shards, shardSize)
	if err != nil {
		return
	}

	for i := range fhs {
		fhs[i].Seek(0, 0)
	}
	err = c.upload(fhs, meta)

	return
}

func (c *Client) upload(fhs []*os.File, meta metafile.Meta) error {
	shards := len(fhs)
	nodes, err := c.GetNodeClients("")
	if err != nil {
		return err
	}

	shardIdCh := make(chan int, shards)
	for i := 0; i < shards; i++ {
		shardIdCh <- i
	}
	close(shardIdCh)
	errCh := make(chan error, shards)
	wg := sync.WaitGroup{}
	wg.Add(c.BlockUpWorkerCount)
	worker := func() {
		defer wg.Done()
		for id := range shardIdCh {
			retry := 0
		lazyTry:
			nodeIdx := (id + retry) % len(nodes)
			node := nodes[nodeIdx]

			mr := bytes.NewBuffer(meta.Encode(id))
			var r io.Reader
			if retry == 0 {
				r = io.MultiReader(mr, fhs[id])
			} else {
				fhs[id].Seek(0, 0)
				r = io.MultiReader(mr, fhs[id])
			}

			blkHash, err := node.Client.Add(r)
			if err != nil {
				if retry < len(nodes) {
					retry++
					log.Println("err:", err, "retrying... retry times:", retry)
					goto lazyTry
				}
				errCh <- err
			}
			log.Printf("block hash: %s, node id: %s, err: %#v", blkHash, node.Id, err)
		}
	}
	for i := 0; i < c.BlockUpWorkerCount; i++ {
		go worker()
	}
	wg.Wait()
	close(errCh)
	for err1 := range errCh {
		return err1
	}
	return nil
}
