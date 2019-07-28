package rpc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

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

	fhs, err := mgr.RsEncode(r, file.DefaultMaxFsizeInMem)
	if err != nil {
		return
	}

	cid, err = file.GetCidV1(h)
	if err != nil {
		return
	}

	meta := metafile.NewMeta(fname, cid, fsize, uint32(dataShards), uint32(parShards))
	meta.WalletPubKey = c.WalletPubKey
	shardSize += int64(len(meta.Encode(0)))
	shards := dataShards + parShards
	_, err = c.NewUploadJob(cid, fsize, shards, shardSize)
	if err != nil {
		log.Println("NewUploadJob Error:", err)
		return
	}

	err = c.upload(fhs, meta)
	if err != nil {
		err = c.DeleteFile(cid)
	}
	return
}

func (c *Client) upload(fhs []file.File, meta metafile.Meta) error {
	shards := len(fhs)
	nodes, err := c.GetNodeClients("")
	if err != nil {
		return err
	}

	//close all files
	defer func() {
		for i := 0; i < shards; i++ {
			fhs[i].Close()
		}
	}()

	//id/err chanel
	shardIdCh := make(chan int, shards)
	for i := 0; i < shards; i++ {
		shardIdCh <- i
	}
	close(shardIdCh)
	errCh := make(chan error, shards)

	//if failed cancel instantly
	ctx, cancel := context.WithCancel(context.Background())
	rand.Seed(time.Now().UnixNano())
	//upload wokers
	wg := sync.WaitGroup{}
	wg.Add(c.BlockUploadWorkers)
	worker := func() {
		defer wg.Done()
		for id := range shardIdCh {
			select {
			case <-ctx.Done():
				return
			default:
			}
			retry := 0

		lazyTry:
			nodeIdx := rand.Int() % len(nodes)
			node := nodes[nodeIdx]

			mr := bytes.NewBuffer(meta.Encode(id))
			var r io.Reader
			if retry != 0 {
				fhs[id].Seek(0, 0)
			}
			r = io.MultiReader(mr, fhs[id])

			blkHash, err := node.Client.Add(r)
			if err != nil {
				if retry < len(nodes) {
					retry++
					log.Println("err: ", err, "retrying... retry times:", retry)
					goto lazyTry
				}
				errCh <- err
				cancel()
			}
			log.Printf("block hash: %s, node id: %s, upload err: %+v", blkHash, node.Id, err)
		}
	}
	for i := 0; i < c.BlockUploadWorkers; i++ {
		go worker()
	}
	wg.Wait()
	close(errCh)

	if ctx.Err() != nil {
		return <-errCh
	}
	return nil
}
