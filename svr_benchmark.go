package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	dir := flag.String("dir", ".", "files dir")
	workers := flag.Int("workers", 10, "worker counts")
	flag.Parse()
	fmt.Printf("dir: %s, wokers: %d\n", *dir, *workers)

	start := time.Now()
	uploadDir(*dir, *workers)
	fmt.Printf("upload time elapsed: %s", time.Now().Sub(start))
}

func uploadDir(dirPath string, workers int) error {

	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	dir.Close()

	namesChan := make(chan string, len(names))
	for _, name := range names {
		namesChan <- name
	}
	close(namesChan)

	if len(names) < workers {
		workers = len(names)
	}
	resChan := make(chan interface{}, len(names))
	for i := 0; i < workers; i++ {
		go func() {
			for name := range namesChan {
				file, err := os.Open(filepath.Join(dirPath, name))
				if err != nil {
					resChan <- err
					continue
				}

				res, err := upload(file)
				file.Close()
				if err != nil {
					resChan <- err
				} else {
					resChan <- res
				}
			}
		}()
	}

	for i := 0; i < len(names); i++ {
		res := <-resChan
		if uRet, ok := res.(UploadRet); ok {
			fmt.Println("upload Success:", uRet.Cid)
		} else {
			fmt.Println("upload Failed:", res)
		}
	}
	return nil
}

type UploadRet struct {
	Cid string `json:"cid"`
}

func upload(r io.Reader) (ret UploadRet, err error) {
	addr := "http://127.0.0.1:9090/file/upload"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.SetBoundary(strconv.Itoa(time.Now().Nanosecond()))

	part, err := writer.CreateFormFile("file", strconv.Itoa(time.Now().Nanosecond()))
	if err != nil {
		return
	}
	_, err = io.Copy(part, r)
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", addr, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		e := struct {
			Err string `json:"err"`
		}{}
		err = dec.Decode(&e)
		if err != nil {
			return
		}
		err = errors.New(e.Err)
	} else {
		err = dec.Decode(&ret)
	}
	return
}
