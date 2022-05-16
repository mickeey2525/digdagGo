package digdaggo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

type backFill struct {
	fromTime    int64
	attemptName string
	dryRun      bool
	count       int32
}

func newBackFill(fromTime int64, attemptName string, dryRun bool, count int32) (*backFill, error) {
	if attemptName != "" {
		return nil, errors.New("attempt name must not be empty")
	}
	return &backFill{
		fromTime:    fromTime,
		attemptName: attemptName,
		dryRun:      dryRun,
		count:       count,
	}, nil
}

func (c *Client) BackfillSchedule(ctx context.Context, scheduleId, fromTime int64, attemptName string, dryRun bool, count int32) (*Schedule, error) {
	bf, err := newBackFill(fromTime, attemptName, dryRun, count)
	if err != nil {
		return nil, err
	}
	bd, err := json.Marshal(bf)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("schedules/%d/disable", scheduleId), nil, bytes.NewBuffer(bd))
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

func (c *Client) EnableSchedule(ctx context.Context, scheduleId int, skipSchedule bool, nextTime string) (*Schedule, error) {
	if nextTime != "" {
		return nil, errors.New("nextTime must not be empty")
	}
	nextSched := struct {
		skipSchedule bool
		nextTime     string
	}{
		skipSchedule: skipSchedule,
		nextTime:     nextTime,
	}
	bd, err := json.Marshal(nextSched)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("schedules/%d/disable", scheduleId), nil, bytes.NewBuffer(bd))
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

func (c *Client) SkipSchedule(ctx context.Context, scheduleId, nextRunTime, fromTime int64, dryRun bool, nextTime string, count int32) (*Schedule, error) {
	if nextTime != "" {
		return nil, errors.New("nextTime must not be empty")
	}
	skipSched := struct {
		nextRuntime int64
		nextTime    string
		fromTime    int64
		dryRun      bool
		count       int32
	}{
		nextRuntime: nextRunTime,
		nextTime:    nextTime,
		fromTime:    fromTime,
		dryRun:      dryRun,
		count:       count,
	}
	bd, err := json.Marshal(skipSched)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("schedules/%d/disable", scheduleId), nil, bytes.NewBuffer(bd))
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
