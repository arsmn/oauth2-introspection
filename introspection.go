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

func (o *OAuth2Introspection) Introspect(token string) (*OAuth2IntrospectionResult, error) {

	if token == "" {
		return nil, errors.WithStack(ErrMalformedToken)
	}

	var result *OAuth2IntrospectionResult
	var isCacheEnabled = o.c.CacheProvider != nil

	if isCacheEnabled {
		res, err := o.getCache(token)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = res
	}

	if result == nil {
		res, err := o.sendRequest(token)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if isCacheEnabled {
			err := o.setCache(token, result)
			if err != nil {
				return nil, err
			}
		}

		result = res
	}

	return result, nil
}

func (o *OAuth2Introspection) sendRequest(token string) (*OAuth2IntrospectionResult, error) {
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
		if !containsSring(i.Audience, audience) {
			return nil, errors.WithStack(ErrForbidden)
		}
	}

	if len(o.c.Issuers) > 0 {
		if !containsSring(o.c.Issuers, i.Issuer) {
			return nil, errors.WithStack(ErrForbidden)
		}
	}

	return &i, nil
}

func (o *OAuth2Introspection) getCache(token string) (*OAuth2IntrospectionResult, error) {
	val, err := o.c.CacheProvider.Get(token)
	if err != nil {
		return nil, err
	}

	if len(val) == 0 {
		return nil, nil
	}

	var i OAuth2IntrospectionResult
	if err := json.Unmarshal(val, &i); err != nil {
		return nil, err
	}

	return &i, nil
}

func (o *OAuth2Introspection) setCache(token string, result *OAuth2IntrospectionResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return o.c.CacheProvider.Set(token, data)
}
