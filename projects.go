package digdaggo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Revision    string      `json:"revision"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	DeletedAt   interface{} `json:"deletedAt"`
	ArchiveType string      `json:"archiveType"`
	ArchiveMd5  string      `json:"archiveMd5"`
}

type Projects struct {
	Projects []Project `json:"projects"`
}

type Revisions struct {
	Revisions []Revision `json:"revisions"`
}
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	IPAddress string `json:"ip_address"`
}
type Td struct {
	User User `json:"user"`
}
type UserInfo struct {
	Td Td `json:"td"`
}
type Revision struct {
	Revision    string    `json:"revision"`
	CreatedAt   time.Time `json:"createdAt"`
	ArchiveType string    `json:"archiveType"`
	ArchiveMd5  string    `json:"archiveMd5"`
	UserInfo    UserInfo  `json:"userInfo"`
}

func (c *Client) GetProjects(ctx context.Context, projectName string) (*Projects, error) {
	parameters := map[string]string{}
	if projectName != "" {
		parameters["name"] = projectName
	}
	req, err := c.newRequest(ctx, "GET", "projects", parameters, nil, nil)
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
	var projects Projects
	err = c.decodeBody(resp, &projects)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}

func (c *Client) GetProjectsWithID(ctx context.Context, projectId string) (*Project, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s", projectId), nil, nil, nil)
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
	var project Project
	err = c.decodeBody(resp, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) PutProject(ctx context.Context, filepath, projectName string) (*Project, error) {
	if projectName == "" {
		return nil, errors.New("project name is required")
	}

	revision := uuid.New().String()
	parameters := map[string]string{}

	parameters["project"] = projectName
	parameters["revision"] = revision

	digFiles, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func(digFiles *os.File) {
		err := digFiles.Close()
		if err != nil {

		}
	}(digFiles)

	header := map[string]string{"content-type": "application/gzip"}

	req, err := c.newRequest(ctx, "PUT", "projects", parameters, digFiles, header)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}
	var project Project
	err = c.decodeBody(resp, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) DeleteProjectsWithID(ctx context.Context, projectId string) (*Project, error) {
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("projects/%s", projectId), nil, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}
	var project Project
	err = c.decodeBody(resp, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) DownloadProjectFiles(ctx context.Context, projectId, revision, destPath string, directDownload bool) error {
	downloadOption := strconv.FormatBool(directDownload)
	parameters := map[string]string{}

	parameters["revision"] = revision
	parameters["direct_download"] = downloadOption
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/archive", projectId), parameters, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return checkStatus
	}
	fmt.Println("-----Started to download------")
	if destPath == "" {
		destPath, err = filepath.Abs(".")
		if err != nil {
			return err
		}
	}
	er := c.unarchive(destPath, resp.Body)
	if er != nil {
		return er
	}
	fmt.Println("-----Finished-----")
	return nil
}

func (c *Client) GetListRevisions(ctx context.Context, projectId string) (*Revisions, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/revisions", projectId), nil, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}
	var revisions Revisions
	err = c.decodeBody(resp, &revisions)
	if err != nil {
		return nil, err
	}
	return &revisions, nil

}

type ShortProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Workflow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetProjectsSchedules(ctx context.Context, projectId, workflow, lastId string) (*ScheduleList, error) {
	parameters := map[string]string{}
	if workflow != "" {
		parameters["workflow"] = workflow
	}
	if lastId != "" {
		parameters["last_id"] = lastId
	}

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/schedules", projectId), parameters, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}
	var schedules ScheduleList
	err = c.decodeBody(resp, &schedules)
	if err != nil {
		return nil, err
	}
	return &schedules, nil
}

type Secrets struct {
	Secrets []interface{} `json:"secrets"`
}

func (c *Client) GetSecrets(ctx context.Context, projectId string) (*Secrets, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/secrets", projectId), nil, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}
	var secrets Secrets
	err = c.decodeBody(resp, &secrets)
	if err != nil {
		return nil, err
	}
	return &secrets, nil

}

func (c *Client) PutSecrets(ctx context.Context, projectId string, secrets map[string]string) error {
	jsn, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, "PUT", fmt.Sprintf("projects/%s/secrets", projectId), nil, bytes.NewReader(jsn), nil)
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

		}
	}(resp.Body)

	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return checkStatus
	}
	return nil
}

func (c *Client) DeleteSecret(ctx context.Context, projectId, key string) error {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/secrets/%s", projectId, key), nil, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return checkStatus
	}
	return nil
}

// Sessions List of the Sessions
type Sessions struct {
	Sessions []Session `json:"sessions"`
}

// TD workflow customized Parameter
type Params struct {
	TdAttemptInitiatedUserID    int    `json:"_td_.attempt_initiated_user_id"`
	TdAttemptInitiatedUserEmail string `json:"_td_.attempt_initiated_user_email"`
	TdAttemptInitiatedUserIP    string `json:"_td_.attempt_initiated_user_ip"`
	TdRevisionCreatedUserID     int    `json:"_td_.revision_created_user_id"`
	TdRevisionCreatedUserEmail  string `json:"_td_.revision_created_user_email"`
	TdRevisionCreatedUserIP     string `json:"_td_.revision_created_user_ip"`
}

// The Last Attempt which is contained in Session Object
type LastAttempt struct {
	ID               string      `json:"id"`
	RetryAttemptName interface{} `json:"retryAttemptName"`
	Done             bool        `json:"done"`
	Success          bool        `json:"success"`
	CancelRequested  bool        `json:"cancelRequested"`
	Params           Params      `json:"params"`
	CreatedAt        time.Time   `json:"createdAt"`
	FinishedAt       time.Time   `json:"finishedAt"`
}

// This struct embodies Session
type Session struct {
	ID          string       `json:"id"`
	Project     ShortProject `json:"project"`
	Workflow    Workflow     `json:"workflow"`
	SessionUUID string       `json:"sessionUuid"`
	SessionTime time.Time    `json:"sessionTime"`
	LastAttempt LastAttempt  `json:"lastAttempt"`
}

func (c *Client) GetProjectSessions(ctx context.Context, projectId, workflowName, lastId, pageSize string) (*Sessions, error) {
	parameters := map[string]string{}
	if workflowName != "" {
		parameters["workflow"] = workflowName
	}
	if lastId != "" {
		parameters["last_id"] = lastId
	}
	if pageSize != "" {
		parameters["page_size"] = pageSize
	}

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/sessions", projectId), parameters, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}

	var sessions Sessions
	err = c.decodeBody(resp, &sessions)
	if err != nil {
		return nil, err
	}
	return &sessions, nil
}

type Workflows struct {
	Workflows []ProjectWorkflow `json:"workflows"`
}

type Config interface{}

type ProjectWorkflow struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Project  Project `json:"project"`
	Revision string  `json:"revision"`
	Timezone string  `json:"timezone"`
	Config   Config  `json:"config"`
}

func (c *Client) GetProjectWorkflows(ctx context.Context, projectId, revision, workflowName string) (*Workflows, error) {
	parameters := map[string]string{}
	if revision != "" {
		parameters["workflow"] = revision
	}
	if workflowName != "" {
		parameters["name"] = workflowName
	}

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/workflows", projectId), parameters, nil, nil)
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

		}
	}(resp.Body)
	checkStatus := c.checkHttpResponseCode(resp)
	if checkStatus != nil {
		return nil, checkStatus
	}

	var workflows Workflows
	err = c.decodeBody(resp, &workflows)
	if err != nil {
		return nil, err
	}
	return &workflows, nil
}
