package client

import (
	"crypto/tls"
	"net/http"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/nolte/terraform-provider-harbor/gen/harborctl/client"
)

// NewClient creates common settings
func NewClient(host string, username string, password string,
	insecure bool, basepath string, schema string) *apiclient.Harbor {
	basicAuth := httptransport.BasicAuth(username, password)

	// allow skipping ssl
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: insecure, // nolint: gosec
	}

	apiSchemes := make([]string, 1)
	apiSchemes[0] = schema

	transport := httptransport.New(host, basepath, apiSchemes)

	// add default auth
	transport.DefaultAuthentication = basicAuth

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)

	return client
}
