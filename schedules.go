package digdaggo

import (
	"context"
	"fmt"
	"time"
)

type Schedule struct {
	ID      string `json:"id"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Workflow struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"workflow"`
	NextRunTime      time.Time `json:"nextRunTime"`
	NextScheduleTime time.Time `json:"nextScheduleTime"`
	DisabledAt       time.Time `json:"disabledAt"`
}

type ScheduleList struct {
	Schedules []Schedule
}

func (c *Client) GetSchedules(ctx context.Context, lastId string) (*ScheduleList, error) {
	parameters := map[string]string{}
	if lastId != "" {
		parameters["last_id"] = lastId
	}
	req, err := c.newRequest(ctx, "GET", "schedules", parameters, nil)
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
	var schedules ScheduleList
	err = c.decodeBody(resp, &schedules)
	if err != nil {
		return nil, err
	}
	return &schedules, nil
}

func (c *Client) GetScheduleWithId(ctx context.Context, scheduleId int) (*Schedule, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("schedules/%d", scheduleId), nil, nil)
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
	var schedule Schedule
	err = c.decodeBody(resp, &schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (c *Client) DisableScheduleWithId(ctx context.Context, scheduleId int) (*Schedule, error) {
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("schedules/%d/disable", scheduleId), nil, nil)
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
	var schedule Schedule
	err = c.decodeBody(resp, &schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
