package rpc

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/metafile"
	"github.com/ipweb-group/go-sdk/contracts/storage"
	"github.com/ipweb-group/go-sdk/file"
	"github.com/ironsmile/nedomi/utils"
)

func (c *Client) Download(fileHash string, w io.Writer) (metaAll metafile.Meta, err error) {
	blocksInfo, err := c.GetBlocksInfo(fileHash)
	if err != nil {
		return
	}

	meta, err := c.GetMeta(blocksInfo)
	if err != nil {
		return
	}

	dataShards := int(meta.MetaHeader.DataShards)
	parShards := int(meta.MetaHeader.ParShards)
	shards := dataShards + parShards
	dataFhs, hasBroken, err := c.download(blocksInfo[:dataShards], len(meta.Encode(0)))
	if err != nil {
		return
	}
	if hasBroken {
		log.Println("have broken shards, start download parity shards.")
		parFhs, _, e := c.download(blocksInfo[dataShards:], len(meta.Encode(0)))
		if e != nil {
			err = e
			return
		}
		log.Println("parity shards download success.")
		mgr, e := file.NewBlockMgr(dataShards, parShards)
		if err != nil {
			err = e
			return
		}

		fhs := append(dataFhs, parFhs...)
		rdrs := make([]io.Reader, shards)
		wtrs := make([]io.Writer, shards)
		for i := range fhs {
			if fhs[i] == nil {
				rdrs[i] = nil
				wtrs[i], err = file.NewTmpFile()
				if err != nil {
					return
				}
			} else {
				rdrs[i] = fhs[i]
			}
		}

		log.Println("reconstruct start.")
		err = mgr.StreamEncoder.Reconstruct(rdrs, wtrs)
		if err != nil {
			return
		}
		log.Println("reconstruct end.")

		for i := range fhs {
			if fhs[i] == nil {
				fhs[i] = wtrs[i].(file.File)
			}
			fhs[i].Seek(0, 0)
		}
		dataFhs = fhs[:dataShards]
	}

	log.Println("download over.")
	drs := make([]io.Reader, dataShards)
	for i := range dataFhs {
		drs[i] = dataFhs[i]
	}
	fr := io.LimitReader(io.MultiReader(drs...), meta.FSize)
	_, err = io.Copy(w, fr)
	return
}

func (c *Client) StreamRead(fileHash string) (rc io.ReadCloser, metaAll metafile.Meta, err error) {
	blocksInfo, err := c.GetBlocksInfo(fileHash)
	if err != nil {
		return
	}

	meta, err := c.GetMeta(blocksInfo)
	if err != nil {
		return
	}

	dataShards := int(meta.MetaHeader.DataShards)
	_, _, shardSize := file.BlockCount(meta.FSize)
	lastShardSize := meta.FSize - int64(dataShards-1)*shardSize
	rcs := make([]io.ReadCloser, dataShards)
	retryTimes := 0
	var partSize int64 = 0
	for i := 0; i < dataShards; i++ {
		blockInfo := blocksInfo[i]
		nodes, e := c.GetNodes(blockInfo.PeerId)
		if e != nil {
			err = e
			log.Printf("GetNode error: %s, node id: %s, block hash: %s \n", blockInfo.PeerId, blockInfo.BlockHash, err)
			return
		}

		if i == dataShards-1 {
			partSize = lastShardSize
		}

		node := nodes[0]
	lazyTry:
		rc1, e := ReadAt(node.Client, blockInfo.BlockHash, int64(meta.Len()), partSize)
		if e != nil {
			err = e
			if retryTimes < 3 {
				log.Printf("read failed, node id: %s, block hash: %s, error: %#v , retryTimes: %d, retrying\n", blockInfo.PeerId, blockInfo.BlockHash, err, retryTimes)
				node = getRandonNode(nodes)
				retryTimes++
				goto lazyTry
			}
			return
		}

		rcs[i] = rc1
		log.Printf("read block, node id: %s, Block Hash: %s, error: %s \n", node.Id, blockInfo.BlockHash, err)
	}

	rc = utils.MultiReadCloser(rcs...)
	return
}

func (c *Client) download(blocksInfo []storage.BlockInfo, metaLen int) (fhs []file.File, hasBroken bool, err error) {
	blockNum := len(blocksInfo)
	fhs, err = file.NewTmpFiles(blockNum)
	if err != nil {
		return
	}

	sem := make(chan bool, c.BlockDownloadWorkers)
	for i := 0; i < blockNum; i++ {
		sem <- true
		go func(idx int) (err error) {
			log.Println("block idx:", idx)
			defer func() {
				<-sem
				if err != nil {
					log.Printf("block download failed idx: %d, error: %#v\n", idx, err)
					fhs[idx].Close()
					fhs[idx] = nil
					hasBroken = true
				} else {
					fhs[idx].Seek(0, 0)
				}
			}()

			blockInfo := blocksInfo[idx]
			node, err := c.GetNode(blockInfo.PeerId)
			if err != nil {
				return
			}

			dialRetryTimes := 0
			downloadRetryTimes := 0
		lazyTry:
			dialStart := time.Now()
			rc1, err := ReadAt(node.Client, blockInfo.BlockHash, int64(metaLen), 0)
			log.Printf("dial to node finished, idx: %d, node id: %s, block hash: %s, time elapsed: %s,  error: %#v \n", idx, node.Id, blockInfo.BlockHash, time.Now().Sub(dialStart), err)
			if err != nil {
				if dialRetryTimes < 2 {
					log.Printf("retrying, idx: %d, block hash: %s, new node id: %s, diaRetryTimes: %d\n", idx, blockInfo.BlockHash, node.Id, dialRetryTimes)
					dialRetryTimes++
					goto lazyTry
				}

				log.Printf("retry failed, return, idx: %d, diaRetryTimes: %d \n", idx, dialRetryTimes)
				return
			}
			if rc1 == nil {
				log.Printf("get nil reader from node, node down, idx: %d, node id: %s, block hash: %s\n", idx, node.Id, blockInfo.BlockHash)
				err = errors.New("ipfs node down")
				return
			}

			copyStart := time.Now()
			_, err = io.Copy(fhs[idx], rc1)
			log.Printf("block download finished, idx: %d, node id: %s, block hash: %s, time elapsed: %s, error: %#v \n", idx, node.Id, blockInfo.BlockHash, time.Now().Sub(copyStart), err)
			if err != nil {
				if downloadRetryTimes < 3 {
					log.Printf("retrying, idx: %d, downloadRetryTimes: %d \n", idx, downloadRetryTimes)
					downloadRetryTimes++
					fhs[idx].Seek(0, 0)
					goto lazyTry
				}
			}
			return
		}(i)
	}
	//wait
	for i := 0; i < c.BlockDownloadWorkers; i++ {
		sem <- true
	}
	return
}

func (c *Client) GetMeta(bis []storage.BlockInfo) (meta *metafile.Meta, err error) {
	for _, bi := range bis {
		nodes, err := c.GetNodes(bi.PeerId)
		if err != nil {
			return nil, err
		}

		for _, n := range nodes {
			meta, err := getMeta(n.Client, bi.BlockHash)
			if err == nil {
				return meta, err
			}
		}
	}
	return
}

func getMeta(node *shell.Shell, blkHash string) (m *metafile.Meta, err error) {
	rc1, err := ReadAt(node, blkHash, 0, metafile.MetaHeaderLength)
	if err != nil {
		return
	}
	defer rc1.Close()
	metaHeaderB, err := ioutil.ReadAll(rc1)
	if err != nil {
		return
	}
	metaHeader, err := metafile.DecodeMetaHeader(metaHeaderB)
	if err != nil {
		return
	}

	metaLen := metafile.MetaHeaderLength + metaHeader.MetaBodyLength
	rc2, err := ReadAt(node, blkHash, 0, int64(metaLen))
	if err != nil {
		return
	}
	defer rc2.Close()
	metaData, err := ioutil.ReadAll(rc2)
	if err != nil {
		return
	}

	return metafile.DecodeMeta(metaData)
}
