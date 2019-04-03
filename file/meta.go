package file

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"syscall"

	"go-sdk/utils/bytes"
)

const MetaHeaderLength = 16

//block data header fixed length
type MetaHeader struct {
	DataShards     uint32 //4 Byte
	ParShards      uint32 //4 Byte
	BlockIdx       uint32 //4 Byte
	MetaBodyLength uint32 //4 Byte
}

//block data header mutual length
type MetaBody struct {
	FName string `json:"n"`
	FSize int64  `json:"s"`
	FHash string `json:"h"`
}

type Meta struct {
	MetaHeader
	MetaBody
}

func NewMeta(fname, fhash string, fsize int64, dataShards, parShards uint32) Meta {
	metaHeader := MetaHeader{DataShards: dataShards, ParShards: parShards}
	metaBody := MetaBody{FName: fname, FSize: fsize, FHash: fhash}
	return Meta{metaHeader, metaBody}
}

func (m Meta) Encode(blkIdx int) (data []byte) {
	metaBody, _ := json.Marshal(m.MetaBody)
	m.BlockIdx = uint32(blkIdx)
	m.MetaBodyLength = uint32(len(metaBody))
	data = EncodeMetaHeader(&m.MetaHeader)
	data = append(data, metaBody...)
	return
}

func (m Meta) Len() int {
	return int(MetaHeaderLength + m.MetaBodyLength)
}

func DecodeMeta(data []byte) (meta *Meta, err error) {
	if len(data) < MetaHeaderLength {
		err = errors.New("data length less than meta header length")
		return
	}

	metaHeader, err := DecodeMetaHeader(data[:MetaHeaderLength])
	if err != nil {
		return
	}

	metaLength := MetaHeaderLength + metaHeader.MetaBodyLength
	if len(data) < int(metaLength) {
		err = errors.New("data length less than meta length")
		return
	}

	metaBody := MetaBody{}
	err = json.Unmarshal(data[MetaHeaderLength:metaLength], &metaBody)
	meta = &Meta{*metaHeader, metaBody}

	return
}

func EncodeMetaHeader(meta *MetaHeader) []byte {
	fh := make([]byte, MetaHeaderLength)
	w := bytes.NewWriter(fh)
	binary.Write(w, binary.LittleEndian, meta)
	return fh
}

func DecodeMetaHeader(metaBytes []byte) (metaHeader *MetaHeader, err error) {
	if len(metaBytes) != MetaHeaderLength {
		err = syscall.EINVAL
		return
	}
	metaHeader = new(MetaHeader)

	r := bytes.NewReader(metaBytes)
	err = binary.Read(r, binary.LittleEndian, metaHeader)
	return
}
