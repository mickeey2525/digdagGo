package digdaggo

import (
	"context"
	"fmt"
)

type WorkflowsList struct {
	Workflows []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Project struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Revision string      `json:"revision"`
		Timezone string      `json:"timezone"`
		Config   interface{} `json:"config"`
	} `json:"workflows"`
}

type DetailedWorkflow struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Revision string      `json:"revision"`
	Timezone string      `json:"timezone"`
	Config   interface{} `json:"config"`
}

func (c *Client) GetWorkflowList(ctx context.Context, last_id, count string) (*WorkflowsList, error) {
	paramters := map[string]string{}
	if last_id != "" {
		paramters["last_id"] = last_id
	}
	if count != "" {
		paramters["count"] = count
	}

	req, err := c.newRequest(ctx, "GET", "workflows", paramters, nil)
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
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("workflows/%s", workflowId), nil, nil)
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
