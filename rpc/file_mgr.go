package rpc

// Remove the given path
func (c *Client) Remove(fHash string) error {
	return c.RemoveByClientKey(c.Client.GetClientKey(), fHash)
}

func (c *Client) RemoveByClientKey(clientKey string, fHash string) error {
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

	err = c.DeleteFileByClientKey(clientKey, fHash)
	if err != nil && err.Error() == ErrContractNotFound.Error() {
		err = ErrContractNotFound
	}
	return err
}
