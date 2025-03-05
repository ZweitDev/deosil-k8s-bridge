package lib

import (
	"encoding/json"
	"errors"
)

type Command struct {
	Action string `json:"action"`
	Filters struct {
		Namespace string `json:"namesapce,omitempty"`
		Label string `json:"label,omitempty"`
	} `json:"filter,omitempty"`
}

func parseCommand(body []byte) (*Command, error) {
	var cmd Command

	if err := json.Unmarshal(body, &cmd); err != nil {
		return nil, err
	}

	// Validate the action
	switch cmd.Action {
	case "GetPods", "GetNodes", "GetDeployments":
		// Valid action
	default:
		return nil, errors.New("invalid action")
	}

	return &cmd, nil
}