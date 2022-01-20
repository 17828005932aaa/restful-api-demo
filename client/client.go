package client

import (
	"fmt"
	"restful-api-demo/apps/host"

	"google.golang.org/grpc"
)

func NewClient(conf *Config) (*Client, error) {
	//拨号连接到指定客户端
	conn, err := grpc.Dial(conf.Addr,grpc.WithInsecure(),grpc.WithPerRPCCredentials(conf.Authentication))
	if err != nil {
		return nil, err
	}

	return &Client{
		conf: conf,
		conn: conn,
	}, nil
}

type Client struct {
	conf *Config
	conn *grpc.ClientConn
}

//new一个app客户端请求实例
func (c *Client) Host() host.ServiceClient {
	return host.NewServiceClient(c.conn)
}
