package p2p

import (
	"github.com/ipfs/go-ipfs/repo"
	"io/ioutil"
	"os"
)

type Repo struct {
	*repo.Mock
}

func (m Repo) SwarmKey() (key []byte, err error) {
	f, err := os.Open("/Users/wf/.ipfs/swarm.key")
	if err != nil {
		return
	}
	key, err = ioutil.ReadAll(f)
	return
}
