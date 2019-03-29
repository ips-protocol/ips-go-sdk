package file

import (
	"encoding/binary"
	"syscall"

	"go-sdk/utils/bytes"
)

const MetaBytes = 16

//block data header fixed length
type Meta struct {
	DataShards   uint32 //4 Byte
	ParShards    uint32 //4 Byte
	BlockIdx     uint32 //4 Byte
	MetaExLength uint32 //4 Byte
}

//block data header mutual length
type MetaEx struct {
	FName string `json:"f_name"`
	FSize int64  `json:"f_size"`
	FHash string `json:"f_hash"`
}

type MetaAll struct {
	Meta
	MetaEx
}

func EncodeMeta(meta *Meta) []byte {

	fh := make([]byte, MetaBytes)
	w := bytes.NewWriter(fh)
	binary.Write(w, binary.LittleEndian, meta)
	return fh
}

func DecodeMeta(metaBytes []byte) (meta *Meta, err error) {

	if len(metaBytes) != MetaBytes {
		err = syscall.EINVAL
		return
	}
	meta = new(Meta)

	r := bytes.NewReader(metaBytes)
	err = binary.Read(r, binary.LittleEndian, meta)
	return
}
