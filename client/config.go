package client

import "context"

type Config struct {
	Addr string
	*Authentication
}

func NewDefaultConfig() *Config {
	return & Config{
		Addr: "127.0.0.1:3570",
		Authentication: &Authentication{},
	}
}

const (
	ClientHeaderKey = "client-id"
	ClientSecretKey = "client-secret"
)

// Authentication todo
type Authentication struct {
	clientID     string
	clientSecret string
}

// SetClientCredentials todo
func (a *Authentication) SetClientCredentials(clientID, clientSecret string) {
	a.clientID = clientID
	a.clientSecret = clientSecret
}

//设置auth请求鉴权
// GetRequestMetadata todo
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{
		ClientHeaderKey: a.clientID,
		ClientSecretKey: a.clientSecret,
	}, nil
}

// RequireTransportSecurity todo
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
