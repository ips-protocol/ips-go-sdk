package redis

import (
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	conf.LoadConfig("../../websvr/conf.json")
}

func TestGetClient(t *testing.T) {
	client := GetClient()
	fmt.Println(client)

	pong, err := client.Ping().Result()
	assert.NoError(t, err)
	assert.Equal(t, pong, "PONG")
	fmt.Println(pong, err)
}