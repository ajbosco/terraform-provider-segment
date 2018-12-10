package segment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ListDestinations returns all destinations for a source
func (c *Client) ListDestinations(src string) (Destinations, error) {
	var d Destinations
	data, err := c.doRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, src, DestinationEndpoint),
		nil)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return d, errors.Wrap(err, "failed to unmarshal destinations response")
	}

	return d, nil
}

// GetDestination returns information about a destination for a source
func (c *Client) GetDestination(src string, dest string) (Destination, error) {
	var d Destination
	data, err := c.doRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, src, DestinationEndpoint, dest),
		nil)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return d, errors.Wrap(err, "failed to unmarshal destination response")
	}

	return d, nil
}

// CreateDestination creates a new destination for a source
func (c *Client) CreateDestination(src string, dest Destination) (Destination, error) {
	var d Destination
	req := destinationCreateRequest{dest}
	data, err := c.doRequest(http.MethodPost,
		fmt.Sprintf("%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, src, DestinationEndpoint),
		req)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return d, errors.Wrap(err, "failed to unmarshal destination response")
	}

	return d, nil
}

// DeleteDestinaton deletes a destination for a source from the workspace
func (c *Client) DeleteDestinaton(src string, dest string) error {
	_, err := c.doRequest(http.MethodDelete,
		fmt.Sprintf("%s/%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, src, DestinationEndpoint, dest),
		nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDestination updates an existing destination with a new config
func (c *Client) UpdateDestination(src string, dest string, destConfig Destination, updateMask UpdateMask) (Destination, error) {
	var d Destination
	req := destinationUpdateRequest{Destination: destConfig, UpdateMask: updateMask}
	data, err := c.doRequest(http.MethodPatch,
		fmt.Sprintf("%s/%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, src, DestinationEndpoint, dest),
		req)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return d, errors.Wrap(err, "failed to unmarshal destination response")
	}

	return d, nil
}
