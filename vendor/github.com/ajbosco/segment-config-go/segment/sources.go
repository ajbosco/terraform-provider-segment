package segment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ListSources returns all sources for a workspace
func (c *Client) ListSources() (Sources, error) {
	var s Sources
	data, err := c.doRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/%s", WorkspacesEndpoint, c.workspace, SourceEndpoint),
		nil)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return s, errors.Wrap(err, "failed to unmarshal sources response")
	}

	return s, nil
}

// GetSource returns information about a source
func (c *Client) GetSource(src string) (Source, error) {
	var s Source
	data, err := c.doRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, src),
		nil)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return s, errors.Wrap(err, "failed to unmarshal source response")
	}

	return s, nil
}

// CreateSource creates a new source
func (c *Client) CreateSource(src Source) (Source, error) {
	var s Source
	req := sourceCreateRequest{src}
	data, err := c.doRequest(http.MethodPost,
		fmt.Sprintf("%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint),
		req)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(data, &s)
	if err != nil {
		return s, errors.Wrap(err, "failed to unmarshal source response")
	}

	return s, nil
}

// DeleteSource deletes a source from the workspace
func (c *Client) DeleteSource(src string) error {
	_, err := c.doRequest(http.MethodDelete,
		fmt.Sprintf("%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, src),
		nil)
	if err != nil {
		return err
	}

	return nil
}
