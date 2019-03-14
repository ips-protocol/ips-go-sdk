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

type Config struct {
	DataShards int `json:"data_shards"`
	ParShards  int `json:"par_shards"`
}

type BlockMgr struct {
	Config
	reedsolomon.StreamEncoder
}

func NewBlockMgr(cfg Config, o ...reedsolomon.Option) (mgr *BlockMgr, err error) {

	dataShards := cfg.DataShards
	if dataShards > 257 {
		err = errors.New("too many data shards")
		return
	}

	e, err := reedsolomon.NewStream(dataShards, cfg.ParShards, o...)
	if err != nil {
		return
	}

	mgr = &BlockMgr{
		Config:        cfg,
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

func (m *BlockMgr) SplitFile2(fname string) (rcs []io.Reader, err error) {

	f, err := os.Open(fname)
	if err != nil {
		return
	}

	fstat, err := f.Stat()
	if err != nil {
		return
	}
	//defer f.Close()

	fsize := fstat.Size()
	shards := m.DataShards + m.ParShards

	perShard := (fsize + int64(m.DataShards) - 1) / int64(m.DataShards)
	padding := make([]byte, (int64(shards)*perShard)-fsize)
	data := io.MultiReader(f, bytes.NewBuffer(padding))

	rs := make([]io.Reader, m.DataShards)
	parWs := make([]io.Writer, m.ParShards)
	rcs = make([]io.Reader, shards)
	for i := range rcs {
		buf := bytes.NewBuffer(nil)
		if i < m.DataShards {
			r := io.LimitReader(data, perShard)
			rcs[i] = io.TeeReader(r, buf)
			rs[i] = buf
		} else {
			parWs[i-m.DataShards] = buf
			rcs[i] = buf
		}
	}

	err = m.Encode(rs, parWs)
	if err != nil {
		return
	}

	return
}

//func (r BlockMgr) Split(data io.Reader, dst []io.Writer, size int64) error {
//
//	return nil
//}

//func (m *BlockMgr) ConcatFile(fname string) (rcs []io.ReadCloser, err error) {
//
//}
