package client_test

import (
	"context"
	"fmt"
	"restful-api-demo/apps/host"
	"restful-api-demo/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost(t *testing.T) {
	should := assert.New(t)

	client, err := client.NewClient(client.NewDefaultConfig())
	should.NoError(err)

	set, err := client.Host().QueryHost(context.Background(), host.NewQueryHostRequest())
	should.NoError(err)

	fmt.Println(set)
}
