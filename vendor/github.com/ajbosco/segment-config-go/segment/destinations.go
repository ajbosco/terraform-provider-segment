package segment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ListDestinations returns all destinations for a source
func (c *Client) ListDestinations(srcName string) (Destinations, error) {
	var d Destinations
	data, err := c.doRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, srcName, DestinationEndpoint),
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
func (c *Client) GetDestination(srcName string, destName string) (Destination, error) {
	var d Destination
	data, err := c.doRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, srcName, DestinationEndpoint, destName),
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
func (c *Client) CreateDestination(srcName string, destName string, connMode string, enabled bool, configs []DestinationConfig) (Destination, error) {
	var d Destination
	destFullName := fmt.Sprintf("%s/%s/%s/%s/%s/%s",
		WorkspacesEndpoint, c.workspace, SourceEndpoint, srcName, DestinationEndpoint, destName)
	dest := Destination{
		Name:           destFullName,
		ConnectionMode: connMode,
		Enabled:        enabled,
		Configs:        configs,
	}
	req := destinationCreateRequest{dest}
	data, err := c.doRequest(http.MethodPost,
		fmt.Sprintf("%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, srcName, DestinationEndpoint),
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

// DeleteDestination deletes a destination for a source from the workspace
func (c *Client) DeleteDestination(srcName string, destName string) error {
	_, err := c.doRequest(http.MethodDelete,
		fmt.Sprintf("%s/%s/%s/%s/%s/%s",
			WorkspacesEndpoint, c.workspace, SourceEndpoint, srcName, DestinationEndpoint, destName),
		nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDestination updates an existing destination with a new config
func (c *Client) UpdateDestination(srcName string, destName string, enabled bool, configs []DestinationConfig) (Destination, error) {
	var d Destination
	destFullName := fmt.Sprintf("%s/%s/%s/%s/%s/%s",
		WorkspacesEndpoint, c.workspace, SourceEndpoint, srcName, DestinationEndpoint, destName)
	dest := Destination{
		Name:    destFullName,
		Enabled: enabled,
		Configs: configs,
	}
	req := destinationUpdateRequest{dest, UpdateMask{Paths: []string{"destination.config", "destination.enabled"}}}
	data, err := c.doRequest(http.MethodPatch, destFullName, req)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return d, errors.Wrap(err, "failed to unmarshal destination response")
	}

	return d, nil
}
