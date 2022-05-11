package digdaggo

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type Files struct {
	File []struct {
		FileName string    `json:"fileName"`
		FileSize int       `json:"fileSize"`
		TaskName string    `json:"taskName"`
		FileTime time.Time `json:"fileTime"`
		AgentID  string    `json:"agentId"`
		Direct   string    `json:"direct"`
	} `json:"files"`
}

func (c *Client) getLogList(ctx context.Context, attemptId int, task string, direct bool) (*Files, error) {
	parameters := map[string]string{}
	if task != "" {
		parameters["task"] = task
	}
	parameters["direct_download"] = strconv.FormatBool(direct)
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("logs/%d/files", attemptId), parameters, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	checkError := c.checkHttpResponseCode(resp)
	if checkError != nil {
		return nil, checkError
	}
	var files Files
	err = c.decodeBody(resp, &files)
	if err != nil {
		return nil, err
	}
	return &files, nil
}
