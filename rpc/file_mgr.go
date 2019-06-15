package rpc

import (
	"context"

	"github.com/ipfs/go-ipfs-api"
)

// Remove the given path
func FilesRm(s *shell.Shell, path string, recursive, force bool) error {
	return s.Request("files/rm", path).
		Option("recursive", recursive).
		Option("force", force).
		Exec(context.Background(), nil)
}

func (c *Client) Remove(fHash string) error {
	blocksInfo, err := c.GetBlocksInfo(fHash)
	if err != nil {
		return err
	}

	for _, bi := range blocksInfo {
		node, err := c.GetNodeClient(bi.PeerId)
		if err != nil {
			return err
		}

		err = node.Client.Unpin(bi.BlockHash)
		if err != nil {
			return err
		}
	}

	err = c.DeleteFile(fHash)
	if err != nil && err.Error() == ErrContractNotFound.Error() {
		err = ErrContractNotFound
	}
	return err
}
