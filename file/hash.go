package file

import (
	"crypto/sha256"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"io"
)

func GetCID(r io.Reader) (id string, err error) {
	hs, err := hashSha256(r)
	if err != nil {
		return
	}

	ehs, err := multihash.Encode(hs, 0x12)
	if err != nil {
		return
	}

	id = cid.NewCidV0(ehs).String()
	return
}

func hashSha256(r io.Reader) ([]byte, error) {
	h := sha256.New()
	_, err := io.Copy(h, r)
	a := h.Sum(nil)
	return a[0:32], err
}
