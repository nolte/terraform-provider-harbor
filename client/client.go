package client

import (
	"crypto/tls"
	"net/http"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
)

type Client struct {
	url      string
	username string
	password string
	insecure bool
	Client   *apiclient.Harbor
}

// NewClient creates common settings
func NewClient(url string, username string, password string, insecure bool, basepath string) *Client {
	basicAuth := httptransport.BasicAuth(username, password)
	// create the transport
	//proxyTLSClientConfig := &tls.Config{InsecureSkipVerify: true}

	if insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	transport := httptransport.New(url, basepath, apiclient.DefaultSchemes)

	// add default auth
	transport.DefaultAuthentication = basicAuth

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)
	return &Client{
		url:      url,
		username: username,
		password: password,
		insecure: insecure,
		Client:   client,
	}
}
