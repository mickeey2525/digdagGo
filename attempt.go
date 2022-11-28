package digdaggo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	req, err := c.newRequest(ctx, "GET", "attempts", param, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	fmt.Println(resp)
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

type Mode int

const (
	FAILED Mode = iota
	FROM
)

type resume struct {
	AttemptId interface{} `json:"attemptId"`
	Mode      Mode        `json:"mode"`
}

type RetryAttemptBody struct {
	SessionTime      time.Time   `json:"sessionTime"`
	WorkflowId       int64       `json:"workflowId"`
	Resume           resume      `json:"resume"`
	RetryAttemptName string      `json:"retryAttemptName"`
	Params           interface{} `json:"params"`
}

type AttemptBody struct {
	SessionTime time.Time   `json:"sessionTime"`
	WorkflowId  int64       `json:"workflowId"`
	Params      interface{} `json:"params"`
}

func (c *Client) StartAttempt(ctx context.Context, params interface{}, workflowID int64, sessionTime time.Time) (*Attempt, error) {
	startAttemptBody := AttemptBody{
		SessionTime: sessionTime,
		WorkflowId:  workflowID,
		Params:      params,
	}
	jsn, err := json.Marshal(startAttemptBody)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", jsn)
	header := map[string]string{"content-type": "application/json"}
	req, err := c.newRequest(ctx, "PUT", "attempts", nil, bytes.NewBuffer(jsn), header)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
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

func (c *Client) RetryAttempt(ctx context.Context, mode Mode, params interface{}, workflowId int64, attemptId interface{}, sessionTime time.Time) (*Attempt, error) {
	attemptNameUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	if attemptId != nil {
		attemptId = attemptId.(int64)
	}
	if attemptId == nil {
		attemptId = ""
	}
	attemptNameUUIDString := attemptNameUUID.String()
	resumeBody := RetryAttemptBody{
		SessionTime:      sessionTime,
		WorkflowId:       workflowId,
		Params:           params,
		Resume:           resume{AttemptId: attemptId, Mode: mode},
		RetryAttemptName: attemptNameUUIDString,
	}
	jsn, err := json.Marshal(resumeBody)
	if err != nil {
		return nil, err
	}
	header := map[string]string{"content-type": "application/json"}
	fmt.Printf("%s\n", jsn)
	req, err := c.newRequest(ctx, "PUT", "attempts", nil, bytes.NewBuffer(jsn), header)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
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

func (c *Client) GetAttempt(ctx context.Context, attemptId string) (*Attempt, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("attempts/%s", attemptId), nil, nil, nil)
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
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("attempts/%s/kill", attemptId), nil, nil, nil)
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
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("attempts/%s/retries", attemptId), nil, nil, nil)
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
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("attempts/%s/tasks", attemptId), nil, nil, nil)
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
