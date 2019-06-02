package file

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/klauspost/reedsolomon"
)

var ErrShortData = errors.New("short data")

const DefaultBlockSize = 1 << 26 //64MB default
const MaxBlockCount = 170

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

func (m *BlockMgr) Split(data io.Reader, size int64) (fhs []*os.File, err error) {
	if size == 0 {
		return fhs, ErrShortData
	}

	fhs, err = CreateTmpFiles(m.DataShards)
	if err != nil {
		return fhs, err
	}

	ws := make([]io.Writer, m.DataShards)
	for i := range ws {
		ws[i] = fhs[i]
	}

	shards := m.DataShards + m.ParShards
	perShard := (size + int64(m.DataShards) - 1) / int64(m.DataShards)

	padding := make([]byte, (int64(shards)*perShard)-size)
	data = io.MultiReader(data, bytes.NewBuffer(padding))

	for i := range fhs {
		n, err := io.CopyN(fhs[i], data, perShard)
		if err != io.EOF && err != nil {
			return fhs, err
		}
		if n != perShard {
			return fhs, ErrShortData
		}

		_, err = fhs[i].Seek(0, 0)
		if err != nil {
			return fhs, err
		}
	}

	return
}

func DeleteTempFiles(fhs []*os.File) error {
	for i := range fhs {
		fpath := filepath.Join(os.TempDir(), fhs[i].Name())
		err := os.Remove(fpath)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateTmpFiles(count int) (fhs []*os.File, err error) {
	tmpFname := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10))
	fhs = make([]*os.File, count)
	for i := range fhs {
		fname := tmpFname + "." + strconv.Itoa(i)
		fhs[i], err = os.Create(fname)
		if err != nil {
			return
		}
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
			dataBuf := &bytes.Buffer{}
			io.Copy(buf, io.TeeReader(r, dataBuf))
			rs[i] = dataBuf
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

func BlockCount(fsize int64) (dataShards, parShards int, shardSize int64) {
	dataShards = int((fsize + int64(DefaultBlockSize) - 1) / int64(DefaultBlockSize))
	if dataShards > MaxBlockCount {
		dataShards = MaxBlockCount
	}
	parShards = int((dataShards + 2) / 3)

	shardSize = (fsize + int64(dataShards) - 1) / int64(dataShards)
	return
}
