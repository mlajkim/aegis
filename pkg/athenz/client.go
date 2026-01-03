package athenz

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Args struct {
	ZmsURL   string // Base URL of the Athenz ZMS server i.e) https://athenz-zms-server:4443/zms/v1
	CertPath string // Path to client certificate file to claim identity i.e) /var/run/athenz/service.crt
	KeyPath  string // Path to client private key file i.e) /var/run/athenz/service.key
	UserTld  string // TopLevelDomain for user members, if TLD is "user.", then give "user" only
}

type AthenzClient struct {
	httpClient *http.Client
	*Args
}

// Initialize a new Athenz Client with mTLS configuration
func New(c Args) (*AthenzClient, error) {
	cert, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certs: %w", err)
	}

	return &AthenzClient{
		httpClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Disabled as we are using self-signed Athenz root CA certs
				Certificates:       []tls.Certificate{cert},
			},
		}},
		Args: &c,
	}, nil
}

func (c *AthenzClient) Get(endpoint string, params url.Values) (*http.Response, error) {
	u, err := url.Parse(c.ZmsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	u = u.JoinPath(endpoint)

	// Set query parameters if provided, in a clean way (to avoid manual string concatenation):
	if params != nil {
		u.RawQuery = params.Encode()
	}

	return c.httpClient.Get(u.String())
}

func (c *AthenzClient) Post(endpoint string, body interface{}) (*http.Response, error) {
	u, err := url.Parse(c.ZmsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	u = u.JoinPath(endpoint)

	var jsonBytes []byte
	if body != nil {
		jsonBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	return c.httpClient.Post(u.String(), "application/json", bytes.NewBuffer(jsonBytes))
}

func (c *AthenzClient) Delete(endpoint string) (*http.Response, error) {
	u, err := url.Parse(c.ZmsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	u = u.JoinPath(endpoint)

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create delete request: %w", err)
	}

	return c.httpClient.Do(req)
}

func (c *AthenzClient) Put(endpoint string, body interface{}) (*http.Response, error) {
	u, err := url.Parse(c.ZmsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	u = u.JoinPath(endpoint)

	var jsonBytes []byte
	if body != nil {
		jsonBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create put request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}
