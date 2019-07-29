package putPolicy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	LoadAppClients("../websvr/app-clients.json")
}

/**
 * 测试编码生成上传 Token
 */
func TestEncodePutPolicy(t *testing.T) {
	appClient, err := GetClientByAccessKey("lfyMRgbefeeFPxbwAgFJyKaNXLQtURnv")
	if err != nil {
		panic(err)
	}

	// 1. Put policy content
	policy := PutPolicy{
		Deadline:     1563295798,
		CallbackUrl:  "http://localhost:8081",
		CallbackBody: "name=$(fname)&size=$(fsize)&hash=$(hash)&width=$(imageWidth)&height=$(imageHeight)",
	}

	result := appClient.MakePolicyWithPutPolicy(policy)

	assert.Equal(t, "lfyMRgbefeeFPxbwAgFJyKaNXLQtURnv:0584281c3990c9c8f056f2a8576a8fdba3d708d4:eyJkZWFkbGluZSI6MTU2MzI5NTc5OCwiY2FsbGJhY2tVcmwiOiJodHRwOi8vbG9jYWxob3N0OjgwODEiLCJjYWxsYmFja0JvZHkiOiJuYW1lPSQoZm5hbWUpXHUwMDI2c2l6ZT0kKGZzaXplKVx1MDAyNmhhc2g9JChoYXNoKVx1MDAyNndpZHRoPSQoaW1hZ2VXaWR0aClcdTAwMjZoZWlnaHQ9JChpbWFnZUhlaWdodCkifQ==", result)
}

/**
 * 测试解码上传 Token
 */
func TestDecodePutPolicy(t *testing.T) {
	policyString := "lfyMRgbefeeFPxbwAgFJyKaNXLQtURnv:0584281c3990c9c8f056f2a8576a8fdba3d708d4:eyJkZWFkbGluZSI6MTU2MzI5NTc5OCwiY2FsbGJhY2tVcmwiOiJodHRwOi8vbG9jYWxob3N0OjgwODEiLCJjYWxsYmFja0JvZHkiOiJuYW1lPSQoZm5hbWUpXHUwMDI2c2l6ZT0kKGZzaXplKVx1MDAyNmhhc2g9JChoYXNoKVx1MDAyNndpZHRoPSQoaW1hZ2VXaWR0aClcdTAwMjZoZWlnaHQ9JChpbWFnZUhlaWdodCkifQ=="
	decodedPutPolicy, err := DecodePutPolicyString(policyString)

	assert.NoError(t, err)
	assert.Equal(t, decodedPutPolicy.PutPolicy.CallbackUrl, "http://localhost:8081")
	assert.Equal(t, decodedPutPolicy.PutPolicy.CallbackBody, "name=$(fname)&size=$(fsize)&hash=$(hash)&width=$(imageWidth)&height=$(imageHeight)")
}
