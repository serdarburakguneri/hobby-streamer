package auth

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type jwksCache struct {
	url        string
	httpClient *http.Client
	keys       map[string]*rsa.PublicKey
	mu         sync.RWMutex
	lastFetch  time.Time
}

func newJWKSCache(url string) *jwksCache {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	return &jwksCache{
		url:        url,
		httpClient: &http.Client{Transport: tr, Timeout: 30 * time.Second},
		keys:       make(map[string]*rsa.PublicKey),
	}
}

func (c *jwksCache) getKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	if key, ok := c.keys[kid]; ok && time.Since(c.lastFetch) < 5*time.Minute {
		c.mu.RUnlock()
		return key, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, errors.NewInternalError("failed to create JWKS request", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch JWKS", err)
	}
	defer resp.Body.Close()

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, errors.NewInternalError("failed to decode JWKS", err)
	}

	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			key, err := convertJWKToPublicKey(jwk)
			if err != nil {
				return nil, errors.NewInternalError("failed to convert JWK to key", err)
			}
			c.keys[kid] = key
			c.lastFetch = time.Now()
			return key, nil
		}
	}

	return nil, errors.NewInternalError(fmt.Sprintf("key %s not found in JWKS", kid), nil)
}

func convertJWKToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)
	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}
