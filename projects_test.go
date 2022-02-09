package digdagGo

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestClient_GetProjects(t *testing.T) {

	testCreatedAt, _ := time.Parse(time.RFC3339, "2018-12-25T15:43:03Z")
	testUpdatedAt, _ := time.Parse(time.RFC3339, "2018-12-25T15:58:52Z")
	tt := []struct {
		name                string
		expectedMethod      string
		expectedRequestPath string
		expectedProjects    *Projects
	}{
		{
			name: "success",

			expectedMethod:      "GET",
			expectedRequestPath: "/projects",
			expectedProjects: &Projects{
				[]Project{
					{
						ID:          "12345",
						Name:        "test",
						Revision:    "c2cdc7e2141f411ba6c6f433d0ab0adf",
						CreatedAt:   testCreatedAt,
						UpdatedAt:   testUpdatedAt,
						DeletedAt:   nil,
						ArchiveType: "s3",
						ArchiveMd5:  "ylv+njP81Shej2RkSHkqkA==",
					},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client, teardown := setup(t, tc.expectedProjects, tc.expectedMethod, tc.expectedRequestPath)
			defer teardown()

			projects, err := client.GetProjects(context.Background(), "")
			if err != nil {
				panic(err)
			}
			for i, v := range projects.Projects {
				if v.ID != tc.expectedProjects.Projects[i].ID ||
					v.Name != tc.expectedProjects.Projects[i].Name ||
					v.Revision != tc.expectedProjects.Projects[i].Revision ||
					v.ArchiveMd5 != tc.expectedProjects.Projects[i].ArchiveMd5 ||
					v.CreatedAt != tc.expectedProjects.Projects[i].CreatedAt ||
					v.UpdatedAt != tc.expectedProjects.Projects[i].UpdatedAt ||
					v.ArchiveType != tc.expectedProjects.Projects[i].ArchiveType ||
					v.DeletedAt != tc.expectedProjects.Projects[i].DeletedAt {
					log.Fatalf("response items wrong. want=%+v, got=%+v", tc.expectedProjects, projects)
				}
			}

		})
	}
}

func TestClient_GetProjectsWithID(t *testing.T) {
	testCreatedAt, _ := time.Parse(time.RFC3339, "2018-12-25T15:43:03Z")
	testUpdatedAt, _ := time.Parse(time.RFC3339, "2018-12-25T15:58:52Z")
	tt := []struct {
		name                string
		expectedMethod      string
		expectedRequestPath string
		expectedProject     *Project
	}{
		{
			name: "success",

			expectedMethod:      "GET",
			expectedRequestPath: "/projects/1",
			expectedProject: &Project{
				ID:          "12345",
				Name:        "test",
				Revision:    "c2cdc7e2141f411ba6c6f433d0ab0adf",
				CreatedAt:   testCreatedAt,
				UpdatedAt:   testUpdatedAt,
				DeletedAt:   nil,
				ArchiveType: "s3",
				ArchiveMd5:  "ylv+njP81Shej2RkSHkqkA==",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client, teardown := setup(t, tc.expectedProject, tc.expectedMethod, tc.expectedRequestPath)
			defer teardown()

			project, err := client.GetProjectsWithID(context.Background(), "1")
			if err != nil {
				panic(err)
			}

			if project.ID != tc.expectedProject.ID ||
				project.Name != tc.expectedProject.Name ||
				project.Revision != tc.expectedProject.Revision ||
				project.ArchiveMd5 != tc.expectedProject.ArchiveMd5 ||
				project.CreatedAt != tc.expectedProject.CreatedAt ||
				project.UpdatedAt != tc.expectedProject.UpdatedAt ||
				project.ArchiveType != tc.expectedProject.ArchiveType ||
				project.DeletedAt != tc.expectedProject.DeletedAt {
				log.Fatalf("response items wrong. want=%+v, got=%+v", tc.expectedProject, project)
			}

		})
	}
}

func setup(t *testing.T, expectedRes interface{}, expectedMethod, expectedRequestPath string) (*Client, func()) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != expectedMethod {
			t.Fatalf("request method wrong. want=%s, got=%s", expectedMethod, req.Method)
		}
		if req.URL.Path != expectedRequestPath {
			t.Fatalf("request path wrong. want=%s, got=%s", expectedRequestPath, req.URL.Path)
		}

		res, err := json.Marshal(expectedRes)
		if err != nil {
			panic(err)
		}
		w.Write(res)
	}))

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to get mock server URL: %s", err.Error())
	}

	cli := &Client{
		BaseURL:    serverURL,
		HTTPClient: server.Client(),
		Logger:     nil,
	}
	teardown := func() {
		server.Close()
	}

	return cli, teardown
}
