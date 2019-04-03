package file

import (
	"bytes"
	"errors"
	"io"

	"github.com/klauspost/reedsolomon"
)

var ErrShortData = errors.New("short data")

const DefaultBlockSize = 1 << 10 //256k for test
const MaxBlockCount = 257

type BlockMgr struct {
	DataShards int
	ParShards  int
	reedsolomon.StreamEncoder
}

func NewBlockMgr(dataShards, parShards int, o ...reedsolomon.Option) (mgr *BlockMgr, err error) {
	if dataShards > MaxBlockCount {
		err = errors.New("too many data shards")
		return
	}

	e, err := reedsolomon.NewStream(dataShards, parShards, o...)
	if err != nil {
		return
	}

	mgr = &BlockMgr{
		DataShards:    dataShards,
		ParShards:     parShards,
		StreamEncoder: e,
	}
	return
}

func (m *BlockMgr) ECShards(reader io.Reader, size int64) (shardsRdr []io.Reader, err error) {
	shards := m.DataShards + m.ParShards
	perShard := (size + int64(m.DataShards) - 1) / int64(m.DataShards)
	padding := make([]byte, (int64(shards)*perShard)-size)
	data := io.MultiReader(reader, bytes.NewBuffer(padding))

	rs := make([]io.Reader, m.DataShards)
	parWs := make([]io.Writer, m.ParShards)
	shardsRdr = make([]io.Reader, shards)
	for i := range shardsRdr {
		buf := &bytes.Buffer{}
		if i < m.DataShards {
			r := io.LimitReader(data, perShard)
			rs[i] = io.TeeReader(r, buf)
			shardsRdr[i] = buf
		} else {
			parWs[i-m.DataShards] = buf
			shardsRdr[i] = buf
		}
	}

	err = m.Encode(rs, parWs)
	if err != nil {
		return
	}

	return
}

func BlockCount(fsize int64) (dataShards, parShards int) {
	dataShards = int((fsize + int64(DefaultBlockSize) - 1) / int64(DefaultBlockSize))
	if dataShards > MaxBlockCount {
		dataShards = MaxBlockCount
	}
	parShards = int((dataShards + 2) / 3)
	return
}