package digdaggo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type ProjectInAttempt struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkflowInAttempt struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Attempt struct {
	Status           string            `json:"status"`
	ID               string            `json:"id"`
	Index            int               `json:"index"`
	Project          ProjectInAttempt  `json:"project"`
	Workflow         WorkflowInAttempt `json:"workflow"`
	SessionID        string            `json:"sessionId"`
	SessionUUID      string            `json:"sessionUuid"`
	SessionTime      time.Time         `json:"sessionTime"`
	RetryAttemptName interface{}       `json:"retryAttemptName"`
	Done             bool              `json:"done"`
	Success          bool              `json:"success"`
	CancelRequested  bool              `json:"cancelRequested"`
	Params           interface{}       `json:"params"`
	CreatedAt        time.Time         `json:"createdAt"`
	FinishedAt       time.Time         `json:"finishedAt"`
}

type AttemptList struct {
	Attempts []Attempt `json:"attempts"`
}

func (c *Client) GetAttempts(ctx context.Context, projectName, workflowName, lastId, pageSize string, includeRetried bool) (*AttemptList, error) {
	param := map[string]string{}

	if projectName != "" {
		param["project"] = projectName
	}

	if workflowName != "" {
		param["workflow"] = workflowName
	}

	if lastId != "" {
		param["lastId"] = lastId
	}

	if pageSize != "" {
		param[pageSize] = pageSize
	}
	if includeRetried {
		param["includeRetried"] = "true"
	}
	req, err := c.newRequest(ctx, "GET", "attempts", param, nil)
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
	var attemptList AttemptList
	err = c.decodeBody(resp, &attemptList)
	if err != nil {
		return nil, err
	}
	return &attemptList, nil
}

type mode int

const (
	FROM mode = iota + 1
	FAILED
)

type resume struct {
	attemptId int64
	mode      mode
}

type attemptBody struct {
	sessionTime      time.Time
	workflowId       int64
	resume           resume
	retryAttemptName string
	params           interface{}
}

func (c *Client) PutAttempt(ctx context.Context, body attemptBody) error {
	jsn, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, "PUT", "attempts", nil, bytes.NewBuffer(jsn))
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return checkStatus
	}
	return nil
}

func (c *Client) GetAttempt(ctx context.Context, attemptId string) (*Attempt, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("attempts/%s", attemptId), nil, nil)
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
	var attempt Attempt
	err = c.decodeBody(resp, &attempt)
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (c *Client) KillAttempt(ctx context.Context, attemptId string) error {
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("attempts/%s/kill", attemptId), nil, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	status := c.checkHttpResponseCode(resp)
	if status != nil {
		return status
	}
	return nil
}

func (c *Client) ListAttempts(ctx context.Context, attemptId string) (*AttemptList, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("attempts/%s/retries", attemptId), nil, nil)
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
	var attemptList AttemptList
	err = c.decodeBody(resp, &attemptList)
	if err != nil {
		return nil, err
	}
	return &attemptList, nil
}

type TasksList struct {
	Tasks []struct {
		ID       string      `json:"id"`
		FullName string      `json:"fullName"`
		ParentID interface{} `json:"parentId"`
		Config   struct {
		} `json:"config"`
		Upstreams       []interface{} `json:"upstreams"`
		State           string        `json:"state"`
		CancelRequested bool          `json:"cancelRequested"`
		ExportParams    struct {
		} `json:"exportParams"`
		StoreParams struct {
		} `json:"storeParams"`
		StateParams struct {
		} `json:"stateParams"`
		UpdatedAt time.Time   `json:"updatedAt"`
		RetryAt   interface{} `json:"retryAt"`
		StartedAt time.Time   `json:"startedAt"`
		Error     struct {
		} `json:"error"`
		IsGroup bool `json:"isGroup"`
	} `json:"tasks"`
}

func (c *Client) ListTasks(ctx context.Context, attemptId string) (*TasksList, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("attempts/%s/tasks", attemptId), nil, nil)
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
	var taskList TasksList
	err = c.decodeBody(resp, &taskList)
	if err != nil {
		return nil, err
	}
	return &taskList, nil
}
