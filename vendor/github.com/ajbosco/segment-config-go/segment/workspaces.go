package segment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// GetWorkspace returns information about a workspace
func (c *Client) GetWorkspace() (Workspace, error) {
	var w Workspace
	data, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%s", WorkspacesEndpoint, c.workspace), nil)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(data, &w)
	if err != nil {
		return w, errors.Wrap(err, "failed to unmarshal workspace response")
	}

	return w, nil
}
