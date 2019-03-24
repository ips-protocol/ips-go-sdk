package main

import (
	//"bytes"
	"encoding/binary"
	"fmt"
	"go-sdk/utils/bytes"
)

func main() {
	// Create a cid manually by specifying the 'prefix' parameters
	//pref := cid.Prefix{
	//	Version:  0,
	//	Codec:    cid.Raw,
	//	MhType:   mh.SHA2_256,
	//	MhLength: -1, // default length
	//}
	//
	//// And then feed it some data
	//c, err := pref.Sum([]byte("Hello World"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Created CID: ", c)

	st := struct {
		A int64
	}{890}

	fh := make([]byte, 8)
	w := bytes.NewWriter(fh)
	binary.Write(w, binary.LittleEndian, &st)

	//s, err := json.Marshal(st)
	fmt.Println("-->", len(fh), fh)

	st1 := &struct {
		A int64
	}{}

	r := bytes.NewReader(fh)

	binary.Read(r, binary.LittleEndian, st1)

	fmt.Printf("====> %d", st1.A)
}
