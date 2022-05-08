package digdaggo

import (
	"context"
	"log"
	"testing"
)

func TestClient_GetWorkflows(t *testing.T) {
	tt := []struct {
		name                 string
		expectedMethod       string
		expectedRequestPath  string
		expectedWorkflowList *WorkflowsList
	}{
		{
			name: "success",

			expectedMethod:      "GET",
			expectedRequestPath: "/workflows",
			expectedWorkflowList: &WorkflowsList{
				[]DetailedWorkflow{
					{
						ID:       "12345",
						Name:     "test",
						Revision: "c2cdc7e2141f411ba6c6f433d0ab0adf",
						Project: ProjectInWorkflow{
							ID:   "12345",
							Name: "test",
						},
						Timezone: "",
					},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client, teardown := setup(t, tc.expectedWorkflowList, tc.expectedMethod, tc.expectedRequestPath)
			defer teardown()

			workflows, err := client.GetWorkflowList(context.Background(), "", "")
			if err != nil {
				panic(err)
			}
			for i, v := range workflows.Workflows {
				if v.ID != tc.expectedWorkflowList.Workflows[i].ID ||
					v.Name != tc.expectedWorkflowList.Workflows[i].Name ||
					v.Revision != tc.expectedWorkflowList.Workflows[i].Revision ||
					v.Project != tc.expectedWorkflowList.Workflows[i].Project ||
					v.Timezone != tc.expectedWorkflowList.Workflows[i].Timezone {
					log.Fatalf("response items wrong. want=%+v, got=%+v", tc.expectedWorkflowList, workflows)
				}
			}

		})
	}
}

func TestClient_GetWorkflowsWithId(t *testing.T) {
	tt := []struct {
		name                 string
		expectedMethod       string
		expectedRequestPath  string
		expectedWorkflowList *DetailedWorkflow
	}{
		{
			name: "success",

			expectedMethod:      "GET",
			expectedRequestPath: "/workflows/12345",
			expectedWorkflowList: &DetailedWorkflow{
				ID:       "12345",
				Name:     "test",
				Revision: "c2cdc7e2141f411ba6c6f433d0ab0adf",
				Project: ProjectInWorkflow{
					ID:   "12345",
					Name: "test",
				},
				Timezone: "",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client, teardown := setup(t, tc.expectedWorkflowList, tc.expectedMethod, tc.expectedRequestPath)
			defer teardown()

			workflow, err := client.GetWorkflowWithID(context.Background(), "12345")
			if err != nil {
				panic(err)
			}
			if workflow.ID != tc.expectedWorkflowList.ID ||
				workflow.Name != tc.expectedWorkflowList.Name ||
				workflow.Revision != tc.expectedWorkflowList.Revision ||
				workflow.Project != tc.expectedWorkflowList.Project ||
				workflow.Timezone != tc.expectedWorkflowList.Timezone {
				log.Fatalf("response items wrong. want=%+v, got=%+v", tc.expectedWorkflowList, workflow)
			}

		})
	}
}
