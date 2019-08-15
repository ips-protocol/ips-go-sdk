package file

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"github.com/klauspost/reedsolomon"
)

var ErrShortData = errors.New("short data")

const DefaultBlockSize = 1 << 26     //64MB default
const DefaultMaxFsizeInMem = 1 << 27 //128MB default
const MaxBlockCount = 170

type BlockMgr struct {
	DataShards int
	ParShards  int
	reedsolomon.StreamEncoder
	reedsolomon.Encoder
}

func NewBlockMgr(dataShards, parShards int, o ...reedsolomon.Option) (mgr *BlockMgr, err error) {
	if dataShards > MaxBlockCount {
		err = errors.New("too many data shards")
		return
	}

	se, err := reedsolomon.NewStream(dataShards, parShards, o...)
	if err != nil {
		return
	}

	be, err := reedsolomon.New(dataShards, parShards, o...)
	if err != nil {
		return
	}

	mgr = &BlockMgr{
		DataShards:    dataShards,
		ParShards:     parShards,
		StreamEncoder: se,
		Encoder:       be,
	}
	return
}

func (m *BlockMgr) RsEncode(r io.Reader, memThreshold int64) (fhs []File, err error) {
	defer func() {
		if err == nil {
			return
		}

		for i := range fhs {
			fhs[i].Close()
		}
	}()

	fh, fsize, err := FileStream(r, memThreshold)
	if err != nil {
		return
	}
	if fsize < memThreshold {
		data, err := ioutil.ReadAll(fh)
		if err != nil {
			return nil, err
		}

		shards, err := m.Encoder.Split(data)
		if err != nil {
			return nil, err
		}

		err = m.Encoder.Encode(shards)
		if err != nil {
			return nil, err
		}

		for _, shard := range shards {
			fhs = append(fhs, newBytesReader(shard))
		}
		return fhs, nil
	}

	dataFhs, err := m.Split(fh, fsize)
	if err != nil {
		return
	}

	parFhs, err := NewTmpFiles(m.ParShards)
	if err != nil {
		return
	}

	dataFhRdrs := make([]io.Reader, m.DataShards)
	for i := range dataFhRdrs {
		dataFhRdrs[i] = dataFhs[i]
	}
	parFhWtrs := make([]io.Writer, m.ParShards)
	for i := range parFhWtrs {
		parFhWtrs[i] = parFhs[i]
	}

	err = m.StreamEncoder.Encode(dataFhRdrs, parFhWtrs)
	if err != nil {
		return
	}
	fhs = append(dataFhs, parFhs...)
	for i := range fhs {
		fhs[i].Seek(0, 0)
	}

	return
}

func (m *BlockMgr) Split(data io.Reader, size int64) (fhs []File, err error) {
	if size == 0 {
		return fhs, ErrShortData
	}
	defer func() {
		if err == nil {
			return
		}

		for i := range fhs {
			fhs[i].Close()
		}
	}()

	fhs, err = NewTmpFiles(m.DataShards)
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

	err = m.StreamEncoder.Encode(rs, parWs)
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
