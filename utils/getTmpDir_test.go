package utils

import (
	"fmt"
	"testing"
)

func TestGetTmpDir(t *testing.T) {
	dir := GetTmpDir()
	fmt.Println(dir)
}
