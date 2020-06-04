package introspection

import "net/http"

// Config holds the configuration for the middleware
type Config struct {
	// Authority is OAuth server address.
	// Required. Default: ""
	Authority string

	// ApiName is name of the API resource used for
	// authentication against introspection endpoint
	// Optional. Default: ""
	APIName string

	// ApiSecret used for authentication against introspection endpoint
	// Optional. Default: ""
	APISecret string

	// Audience defines required audience for authorization.
	// Optional. Default: nil
	Audience []string

	// Issuers defines required issuers for authorization.
	// Optional. Default: nil
	Issuers []string

	// IntrospectionRequestHeaders is list of headers
	// that is send to introspection endpoint.
	// Optional. Default: nil
	IntrospectionRequestHeaders map[string]string
}

type OAuth2Introspection struct {
	c      Config
	client *http.Client
}

type OAuth2IntrospectionResult struct {
	Active    bool                   `json:"active"`
	Extra     map[string]interface{} `json:"ext"`
	Subject   string                 `json:"sub,omitempty"`
	Username  string                 `json:"username"`
	Audience  []string               `json:"aud"`
	TokenType string                 `json:"token_type"`
	Issuer    string                 `json:"iss"`
	ClientID  string                 `json:"client_id,omitempty"`
	Scope     string                 `json:"scope,omitempty"`
}
