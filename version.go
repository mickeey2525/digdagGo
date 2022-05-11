package digdaggo

import (
	"context"
	"errors"
)

type ServerVersion struct {
	Version string `json:"version"`
}

func (c *Client) GetServerVersion(ctx context.Context) (*ServerVersion, error) {
	req, err := c.newRequest(ctx, "GET", "version", nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}
	var serverVersion ServerVersion
	err = c.decodeBody(resp, &serverVersion)
	if err != nil {
		return nil, err
	}
	return &serverVersion, nil
}

type Compatibility struct {
	ServerVersion      string `json:"serverVersion"`
	UpgradeRecommended bool   `json:"upgradeRecommended"`
	APICompatible      bool   `json:"apiCompatible"`
}

func (c *Client) GetServerCompatibility(ctx context.Context, cliVersion string) (*Compatibility, error) {
	param := map[string]string{}
	if cliVersion != "" {
		param["client"] = cliVersion
	} else {
		return nil, errors.New("client version must be specified")
	}
	req, err := c.newRequest(ctx, "GET", "version/check", param, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}
	var serverVersion Compatibility
	err = c.decodeBody(resp, &serverVersion)
	if err != nil {
		return nil, err
	}
	return &serverVersion, nil
}
