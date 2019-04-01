package file

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/klauspost/reedsolomon"
)

var ErrShortData = errors.New("short data")

const DefaultBlockSize = 1 << 10 //256k for test

type BlockMgr struct {
	DataShards int
	ParShards  int
	reedsolomon.StreamEncoder
}

func NewBlockMgr(dataShards, parShards int, o ...reedsolomon.Option) (mgr *BlockMgr, err error) {

	if dataShards > 257 {
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

func (m *BlockMgr) SplitFile(fname string) (rcs []io.ReadCloser, err error) {

	f, err := os.Open(fname)
	if err != nil {
		return
	}

	instat, err := f.Stat()
	if err != nil {
		return
	}
	defer f.Close()

	shards := m.DataShards + m.ParShards
	ws := make([]io.Writer, m.DataShards)
	rs := make([]io.Reader, m.DataShards)
	parWs := make([]io.Writer, m.ParShards)
	rcs = make([]io.ReadCloser, shards)
	for i := range rcs {
		fh, err := ioutil.TempFile(os.TempDir(), fname)
		if err != nil {
			return nil, err
		}

		rcs[i] = fh
		if i < m.DataShards {
			rs[i] = fh
			ws[i] = fh
		} else {
			parWs[i-m.DataShards] = fh
		}
	}

	err = m.Split(f, ws, instat.Size())
	if err != nil {
		panic(err)
	}

	for i := range ws {
		f := ws[i].(*os.File)
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return
		}
	}

	err = m.Encode(rs, parWs)
	if err != nil {
		return
	}

	for i := range rcs {
		f := rcs[i].(*os.File)
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return
		}
	}

	return
}

func (m *BlockMgr) ECShards(reader io.Reader, getMeta func(int) []byte, size int64) (shardsRdr []io.Reader, err error) {

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
			meta := getMeta(i)
			dataWithMeta := io.MultiReader(bytes.NewBuffer(meta), data)
			r := io.LimitReader(dataWithMeta, perShard)
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

func BlockCount(metaSize int, fsize int64, redundancyRate float64) (dataShards, parShards int) {
	blockDataSize := DefaultBlockSize - metaSize
	dataShards = int((fsize + int64(blockDataSize) - 1) / int64(blockDataSize))
	parShards = int(float64(dataShards)*redundancyRate + 1)
	return
}
