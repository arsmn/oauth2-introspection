package introspection

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

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

func New(config ...Config) *OAuth2Introspection {

	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.Authority == "" {
		log.Fatal("Oauth2Introspection: Authority is required")
	}

	client := &http.Client{
		Timeout: time.Millisecond * 500,
	}

	return &OAuth2Introspection{
		c:      cfg,
		client: client,
	}
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

func (o *OAuth2Introspection) Authenticate(token string) (*OAuth2IntrospectionResult, error) {

	if token == "" {
		return nil, errors.WithStack(ErrMalformedToken)
	}
	body := url.Values{"token": {token}}
	introspectReq, err := http.NewRequest(http.MethodPost, o.c.Authority, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for key, value := range o.c.IntrospectionRequestHeaders {
		introspectReq.Header.Set(key, value)
	}

	introspectReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := o.client.Do(introspectReq)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Introspection returned status code %d but expected %d", resp.StatusCode, http.StatusOK)
	}

	var i OAuth2IntrospectionResult
	if err := json.NewDecoder(resp.Body).Decode(&i); err != nil {
		return nil, errors.WithStack(err)
	}

	if len(i.TokenType) > 0 && i.TokenType != "access_token" {
		return nil, errors.WithStack(ErrForbidden)
	}

	if !i.Active {
		return nil, errors.WithStack(ErrUnauthorized)
	}

	for _, audience := range o.c.Audience {
		if !contains(i.Audience, audience) {
			return nil, errors.WithStack(ErrForbidden)
		}
	}

	if len(o.c.Issuers) > 0 {
		if !contains(o.c.Issuers, i.Issuer) {
			return nil, errors.WithStack(ErrForbidden)
		}
	}

	return &i, nil
}

func contains(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}
