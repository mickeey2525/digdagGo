package digdaggo

import (
	"context"
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

func (c *Client) GetAttempts(ctx context.Context, projectName, workflowName, last_id, page_size string, include_retried bool) (*AttemptList, error) {
	paramters := map[string]string{}

	if projectName != "" {
		paramters["project"] = projectName
	}

	if workflowName != "" {
		paramters["workflow"] = workflowName
	}

	if last_id != "" {
		paramters["last_id"] = last_id
	}

	if page_size != "" {
		paramters[page_size] = page_size
	}
	if include_retried {
		paramters["include_retried"] = "true"
	}
	req, err := c.newRequest(ctx, "GET", "attempts", paramters, nil)
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
