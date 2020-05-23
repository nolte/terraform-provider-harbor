package client

import (
	"crypto/tls"
	"net/http"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
)

type Client struct {
	Client *apiclient.Harbor
}

// NewClient creates common settings
func NewClient(host string, username string, password string, insecure bool, basepath string, schema string) *Client {
	basicAuth := httptransport.BasicAuth(username, password)
	// create the transport
	//proxyTLSClientConfig := &tls.Config{InsecureSkipVerify: true}

	if insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	var apiSchemes = []string{schema}
	transport := httptransport.New(host, basepath, apiSchemes)

	// add default auth
	transport.DefaultAuthentication = basicAuth

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)
	return &Client{
		Client: client,
	}
}
