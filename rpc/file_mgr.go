package rpc

import (
	"crypto/sha256"
	"io"

	"github.com/ipweb-group/go-sdk/file"
)

// Remove the given path
func (c *Client) Remove(fHash string) error {
	blocksInfo, err := c.GetBlocksInfo(fHash)
	if err != nil {
		return err
	}

	for _, bi := range blocksInfo {
		node, err := c.GetNode(bi.PeerId)
		if err != nil {
			return err
		}

		err = node.Client.Unpin(bi.BlockHash)
		if err != nil {
			return err
		}
	}

	err = c.DeleteFileByClientKey(c.Client.GetClientKey(), fHash)
	if err != nil && err.Error() == ErrContractNotFound.Error() {
		err = ErrContractNotFound
	}
	return err
}

func (c *Client) GetCid(rdr io.Reader, fname string, fsize int64) (cid string, err error) {
	return c.GetCidByClientKey(c.Client.GetClientKey(), rdr, fname, fsize)
}

func (c *Client) GetCidByClientKey(clientKey string, rdr io.Reader, fname string, fsize int64) (cid string, err error) {
	dataShards, parShards, _ := file.BlockCount(fsize)
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

	_, err = mgr.RsEncode(r, file.DefaultMaxFsizeInMem)
	if err != nil {
		return
	}

	cid, err = file.GetCidV1(h)
	if err != nil {
		return
	}

	return
}
