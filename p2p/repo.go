package p2p

import (
	"fmt"

	"github.com/ipweb-group/go-ipfs/repo"
	"github.com/ipweb-group/go-ipfs/repo/fsrepo"
)

type Repo struct {
	*repo.Mock
}

func (m Repo) SwarmKey() (key []byte, err error) {
	swarm := fsrepo.SwarmKey
	if swarm == "" {
		return nil, fmt.Errorf("not found SwarmKey")
	}
	return []byte(swarm), nil
}
