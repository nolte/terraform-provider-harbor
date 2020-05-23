package client

import (
	"crypto/tls"
	"net/http"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
)

// NewClient creates common settings
func NewClient(host string, username string, password string, insecure bool, basepath string, schema string) *apiclient.Harbor {
	basicAuth := httptransport.BasicAuth(username, password)
	// create the transport
	if insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	apiSchemes := []string{schema}
	transport := httptransport.New(host, basepath, apiSchemes)

	// add default auth
	transport.DefaultAuthentication = basicAuth

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)
	return client
}
