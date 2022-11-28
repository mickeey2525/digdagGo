package digdaggo

import (
	"context"
	"fmt"
)

type WorkflowsList struct {
	Workflows []DetailedWorkflow `json:"workflows"`
}

type ProjectInWorkflow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DetailedWorkflow struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Project  ProjectInWorkflow `json:"project"`
	Revision string            `json:"revision"`
	Timezone string            `json:"timezone"`
	Config   interface{}       `json:"config"`
}

func (c *Client) GetWorkflowList(ctx context.Context, lastId, count string) (*WorkflowsList, error) {
	param := map[string]string{}
	if lastId != "" {
		param["last_id"] = lastId
	}
	if count != "" {
		param["count"] = count
	}

	req, err := c.newRequest(ctx, "GET", "workflows", param, nil, nil)
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
	var workflowsList WorkflowsList
	err = c.decodeBody(resp, &workflowsList)
	if err != nil {
		return nil, err
	}
	return &workflowsList, nil
}

func (c *Client) GetWorkflowWithID(ctx context.Context, workflowId string) (*DetailedWorkflow, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("workflows/%s", workflowId), nil, nil, nil)
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
	var detailedWorkflow DetailedWorkflow
	err = c.decodeBody(resp, &detailedWorkflow)
	if err != nil {
		return nil, err
	}
	return &detailedWorkflow, nil
}
