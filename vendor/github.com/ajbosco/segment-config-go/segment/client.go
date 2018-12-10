package segment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	apiVersion     = "v1beta"
	defaultBaseURL = "https://platform.segmentapis.com"
	mediaType      = "application/json"
)

// Client manages communication with Segment Config API.
type Client struct {
	baseURL     string
	apiVersion  string
	accessToken string
	workspace   string
	client      *http.Client
}

// NewClient creates a new Segment Config API client.
func NewClient(accessToken string, workspace string) *Client {
	return &Client{
		baseURL:     defaultBaseURL,
		apiVersion:  apiVersion,
		accessToken: accessToken,
		workspace:   workspace,
		client:      http.DefaultClient,
	}
}

func (c *Client) doRequest(method, endpoint string, data interface{}) ([]byte, error) {

	// Encode data if we are passed an object.
	b := bytes.NewBuffer(nil)
	if data != nil {
		// Create the encoder.
		enc := json.NewEncoder(b)
		if err := enc.Encode(data); err != nil {
			return nil, errors.Wrap(err, "json encoding data for doRequest failed")
		}
	}

	// Create the request.
	uri := fmt.Sprintf("%s/%s/%s", c.baseURL, c.apiVersion, strings.Trim(endpoint, "/"))
	req, err := http.NewRequest(method, uri, b)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("creating %s request to %s failed", method, uri))
	}

	// Set the proper headers.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", mediaType)

	// Do the request.
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("performing %s request to %s failed", method, uri))
	}
	defer resp.Body.Close()

	// Check that the response status code was OK.
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("invalid access token")
	case http.StatusForbidden:
		return nil, fmt.Errorf("unauthorized access to endpoint")
	case http.StatusNotFound:
		return nil, fmt.Errorf("the requested uri does not exist")
	case http.StatusBadRequest:
		return nil, fmt.Errorf("the request is invalid")
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("decoding response from %s request to %s failed: body -> %s\n", method, uri, string(body)))
	}

	return body, nil
}
