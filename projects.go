package digdagGo

import (
	"context"
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
