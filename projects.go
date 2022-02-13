package digdaggo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	req, err := c.newRequest(ctx, "GET", "projects", parameters, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var projects Projects
		err := c.decodeBody(resp, &projects)
		if err != nil {
			return nil, err
		}
		return &projects, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusForbidden:
		return nil, errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return nil, errors.New("failed to login")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("unexpected error")
	}
}

func (c *Client) GetProjectsWithID(ctx context.Context, projectId string) (*Project, error) {

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s", projectId), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var project Project
		err := c.decodeBody(resp, &project)
		if err != nil {
			return nil, err
		}
		return &project, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusForbidden:
		return nil, errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return nil, errors.New("failed to login")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("unexpected error")
	}
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
	defer digFiles.Close()

	req, err := c.newRequest(ctx, "PUT", "projects", parameters, digFiles)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var project Project
		err := c.decodeBody(resp, &project)
		if err != nil {
			return nil, err
		}
		return &project, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusForbidden:
		return nil, errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return nil, errors.New("failed to login")
	case http.StatusInternalServerError:
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("internal server error. %s", string(bodyBytes))
	default:
		return nil, errors.New("unexpected error")
	}
}

func (c *Client) DeleteProjectsWithID(ctx context.Context, projectId string) (*Project, error) {
	req, err := c.newRequest(ctx, "DELETE", fmt.Sprintf("projects/%s", projectId), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var project Project
		err := c.decodeBody(resp, &project)
		if err != nil {
			return nil, err
		}
		return &project, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusForbidden:
		return nil, errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return nil, errors.New("failed to login")
	case http.StatusNotFound:
		return nil, errors.New("not found the project")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("unexpected error")
	}
}

func (c *Client) DownloadProjectFiles(ctx context.Context, projectId, revision, destPath string, direct_download bool) error {
	download_option := strconv.FormatBool(direct_download)
	parameters := map[string]string{}

	parameters["revision"] = revision
	parameters["direct_donwload"] = download_option
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/archive", projectId), parameters, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
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
	case http.StatusBadRequest:
		return fmt.Errorf("bad Request: %+v", resp.Body)
	case http.StatusForbidden:
		return errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return errors.New("failed to login")
	case http.StatusInternalServerError:
		return errors.New("internal server error")
	default:
		return errors.New("unexpected error")
	}
}

func (c *Client) GetListRevisions(ctx context.Context, projectId string) (*Revisions, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/revisions", projectId), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var revisions Revisions
		err := c.decodeBody(resp, &revisions)
		if err != nil {
			return nil, err
		}
		return &revisions, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusForbidden:
		return nil, errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return nil, errors.New("failed to login")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("unexpected error")
	}
}

type Schedules struct {
	Schedule []Schedule `json:"schedules"`
}
type ShortProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Workflow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Schedule struct {
	ID               string       `json:"id"`
	ShortProject     ShortProject `json:"project"`
	Workflow         Workflow     `json:"workflow"`
	NextRunTime      time.Time    `json:"nextRunTime"`
	NextScheduleTime time.Time    `json:"nextScheduleTime"`
	DisabledAt       interface{}  `json:"disabledAt"`
}

func (c *Client) GetProjectsSchedules(ctx context.Context, projectId, workflow, last_id string) (*Schedules, error) {
	parameters := map[string]string{}
	if workflow != "" {
		parameters["workflow"] = workflow
	}
	if last_id != "" {
		parameters["last_id"] = last_id
	}

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/schedules", projectId), parameters, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var schedules Schedules
		err := c.decodeBody(resp, &schedules)
		if err != nil {
			return nil, err
		}
		return &schedules, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusForbidden:
		return nil, errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return nil, errors.New("failed to login")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	case http.StatusNotFound:
		return nil, fmt.Errorf("not found error %+v", err)
	default:
		return nil, fmt.Errorf("unexpected error: %+v", err)
	}
}

type Secrets struct {
	Secrets []interface{} `json:"secrets"`
}

func (c *Client) GetSecrets(ctx context.Context, projectId string) (*Secrets, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/secrets", projectId), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var secrets Secrets
		err := c.decodeBody(resp, &secrets)
		if err != nil {
			return nil, err
		}
		return &secrets, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusForbidden:
		return nil, errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return nil, errors.New("failed to login")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	case http.StatusNotFound:
		return nil, fmt.Errorf("not found error %+v", err)
	default:
		return nil, fmt.Errorf("unexpected error: %+v", err)
	}
}

func (c *Client) PutSecrets(ctx context.Context, projectId string, secrets map[string]string) error {
	jsn, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	req, err := c.newRequest(ctx, "PUT", fmt.Sprintf("projects/%s/secrets", projectId), nil, bytes.NewReader(jsn))
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return errors.New("bad request")
	case http.StatusForbidden:
		return errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return errors.New("failed to login")
	case http.StatusInternalServerError:
		return errors.New("internal server error")
	case http.StatusNotFound:
		return fmt.Errorf("not found error %+v", err)
	default:
		return fmt.Errorf("unexpected error: %+v", err)
	}
}

func (c *Client) DeleteSecret(ctx context.Context, projectId, key string) error {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("projects/%s/secrets/%s", projectId, key), nil, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return errors.New("bad request")
	case http.StatusForbidden:
		return errors.New("you're not allowed to do this operation")
	case http.StatusUnauthorized:
		return errors.New("failed to login")
	case http.StatusInternalServerError:
		return errors.New("internal server error")
	case http.StatusNotFound:
		return fmt.Errorf("not found error %+v", err)
	default:
		return fmt.Errorf("unexpected error: %+v", err)
	}
}
