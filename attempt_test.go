package digdaggo

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestClient_GetAttempts(t *testing.T) {
	testSessionTime, _ := time.Parse(time.RFC3339, "2018-12-25T15:43:03Z")
	testCreatedAt, _ := time.Parse(time.RFC3339, "2022-04-01T14:00:00:00Z")
	testFinishedAt, _ := time.Parse(time.RFC3339, "2022-04-15T14:00:00:00Z")

	tt := []struct {
		name                string
		expectedMethod      string
		expectedRequestPath string
		expectedAttempts    *AttemptList
	}{
		{
			name: "success",

			expectedMethod:      "GET",
			expectedRequestPath: "/attempts",
			expectedAttempts: &AttemptList{
				[]Attempt{
					{
						Status: "success",
						ID:     "12345",
						Index:  1,
						Workflow: WorkflowInAttempt{
							ID:   "12345",
							Name: "abc123",
						},
						Project: ProjectInAttempt{
							ID:   "12345",
							Name: "test",
						},
						SessionID:        "12345",
						SessionUUID:      "8b7add0f-ea24-420f-a041-135c4c8c4a32",
						SessionTime:      testSessionTime,
						RetryAttemptName: nil,
						Done:             true,
						Success:          true,
						CancelRequested:  false,
						Params:           nil,
						CreatedAt:        testCreatedAt,
						FinishedAt:       testFinishedAt,
					},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client, teardown := setup(t, tc.expectedAttempts, tc.expectedMethod, tc.expectedRequestPath)
			defer teardown()

			attempt, err := client.GetAttempts(context.Background(), "", "", "", "", false)
			if err != nil {
				panic(err)
			}
			for i, v := range attempt.Attempts {
				if v.ID != tc.expectedAttempts.Attempts[i].ID ||
					v.Status != tc.expectedAttempts.Attempts[i].Status ||
					v.Index != tc.expectedAttempts.Attempts[i].Index ||
					v.Project != tc.expectedAttempts.Attempts[i].Project ||
					v.Workflow != tc.expectedAttempts.Attempts[i].Workflow ||
					v.SessionID != tc.expectedAttempts.Attempts[i].SessionID ||
					v.SessionUUID != tc.expectedAttempts.Attempts[i].SessionUUID ||
					v.SessionTime != tc.expectedAttempts.Attempts[i].SessionTime ||
					v.RetryAttemptName != tc.expectedAttempts.Attempts[i].RetryAttemptName ||
					v.Done != tc.expectedAttempts.Attempts[i].Done ||
					v.Success != tc.expectedAttempts.Attempts[i].Success ||
					v.CancelRequested != tc.expectedAttempts.Attempts[i].CancelRequested ||
					v.Params != tc.expectedAttempts.Attempts[i].Params ||
					v.CreatedAt != tc.expectedAttempts.Attempts[i].CreatedAt ||
					v.FinishedAt != tc.expectedAttempts.Attempts[i].FinishedAt {
					log.Fatalf("response items wrong. want=%+v, got=%+v", tc.expectedAttempts, attempt)
				}
			}
		})
	}
}
