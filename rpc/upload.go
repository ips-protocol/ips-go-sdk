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
	return c.UploadWithPathByClientKey(c.Client.GetClientKey(), fpath)
}

func (c *Client) UploadWithPathByClientKey(clientKey string, fpath string) (cid string, err error) {
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

	cid, err = c.uploadByClientKey(clientKey, fh, fname, fi.Size())
	return
}

func (c *Client) Upload(rdr io.Reader, fname string, fsize int64) (cid string, err error) {
	return c.uploadByClientKey(c.Client.GetClientKey(), rdr, fname, fsize)
}

func (c *Client) uploadByClientKey(clientKey string, rdr io.Reader, fname string, fsize int64) (cid string, err error) {
	dataShards, parShards, shardSize := file.BlockCount(fsize)
	mgr, err := file.NewBlockMgr(dataShards, parShards)
	if err != nil {
		return
	}

	h := sha256.New()
	pubKey, err := GetWalletPubKey(clientKey)
	if err != nil {
		return
	}
	_, err = h.Write([]byte(pubKey))
	if err != nil {
		return
	}
	r := io.TeeReader(rdr, h)

	fhs, err := mgr.RsEncode(r, file.DefaultMaxFsizeInMem)
	if err != nil {
		return
	}
	defer func() {
		file.Files(fhs).Close()
	}()

	cid, err = file.GetCidV1(h)
	if err != nil {
		return
	}

	meta := metafile.NewMeta(fname, cid, fsize, uint32(dataShards), uint32(parShards))
	meta.WalletPubKey = pubKey
	shardSize += int64(len(meta.Encode(0)))
	shards := dataShards + parShards
	_, err = c.NewUploadJobByClientKey(clientKey, cid, fsize, shards, shardSize)
	if err != nil {
		log.Println("NewUploadJob Error:", err)
		return
	}

	err = c.upload(fhs, meta)
	if err == nil {
		return
	}

	err1 := c.DeleteFile(cid)
	if err1 != nil {
		err = err1
	}
	return
}

func (c *Client) UploadByClientKey(clientKey string, rdr io.Reader, fname string, fsize int64, preCid string) (cid string, err error) {
	dataShards, parShards, shardSize := file.BlockCount(fsize)
	mgr, err := file.NewBlockMgr(dataShards, parShards)
	if err != nil {
		return
	}

	h := sha256.New()
	pubKey, err := GetWalletPubKey(clientKey)
	if err != nil {
		return
	}
	_, err = h.Write([]byte(pubKey))
	if err != nil {
		return
	}
	r := io.TeeReader(rdr, h)

	fhs, err := mgr.RsEncode(r, file.DefaultMaxFsizeInMem)
	if err != nil {
		return
	}

	cid = preCid

	meta := metafile.NewMeta(fname, cid, fsize, uint32(dataShards), uint32(parShards))
	meta.WalletPubKey = pubKey
	shardSize += int64(len(meta.Encode(0)))
	shards := dataShards + parShards
	_, err = c.NewUploadJobByClientKey(clientKey, cid, fsize, shards, shardSize)
	if err != nil {
		log.Println("NewUploadJob Error:", err)
		return
	}

	err = c.upload(fhs, meta)
	if err == nil {
		return
	}

	err1 := c.DeleteFile(cid)
	if err1 != nil {
		err = err1
	}
	return
}

func (c *Client) upload(fhs []file.File, meta metafile.Meta) error {
	shards := len(fhs)

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
			mr := bytes.NewBuffer(meta.Encode(id))
			var r io.Reader
			if retry != 0 {
				_, err := fhs[id].Seek(0, 0)
				if err != nil {
					return
				}
			}
			r = io.MultiReader(mr, fhs[id])

			_, err := c.Add(r)
			if err != nil {
				if retry < 3 {
					retry++
					log.Println("err: ", err, "retrying... retry times:", retry)
					goto lazyTry
				}
				errCh <- err
				cancel()
			}
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
