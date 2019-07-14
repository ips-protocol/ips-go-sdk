package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	dir := flag.String("dir", ".", "files dir")
	addr := flag.String("addr", "http://127.0.0.1:9090/file/upload", "url")
	fmt.Println("dir:", dir)

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		panic(err)
	}

	start := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(10)
	for _, f := range files {
		go func() {
			defer wg.Done()

			fh, err := os.Open(f.Name())
			if err != nil {
				panic(err)
			}
			req, err := http.NewRequest("POST", *addr, fh)
			if err != nil {
				panic(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			fmt.Println("upload success: ", f.Name(), " : ", string(body))
		}()
	}

	wg.Wait()
	fmt.Println("time elapsed:", time.Now().Sub(start))
}
