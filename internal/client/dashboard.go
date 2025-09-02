package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

const req = "https://lmsug23.iiitkottayam.ac.in/lib/ajax/service.php?sesskey=%v&info=core_calendar_get_action_events_by_timesort"

type Dashboard struct {
	Error bool `json:"error"`
	Data  struct {
		Events []struct {
			Name      string `json:"name"`
			Timestart int    `json:"timestart"`
			Overdue   bool   `json:"overdue"`
			Course    struct {
				Fullname string `json:"fullname"`
			} `json:"course"`
			Purpose string `json:"purpose"`
			URL     string `json:"url"`
		} `json:"events"`
	} `json:"data"`
}

func (lms *LMSCLient) GetDashBoard() (*Dashboard, error) {
	fullReq := fmt.Sprintf(req, lms.Sesskey)
	now := time.Now()
	timesortfrom := now.Add(-24 * time.Hour).Unix()  // 1 day ago
	timesortto := now.Add(7 * 24 * time.Hour).Unix() // 1 week from now

	// JSON payload
	payload := []map[string]any{
		{
			"index":      0,
			"methodname": "core_calendar_get_action_events_by_timesort",
			"args": map[string]any{
				// increase for more limit??
				"limitnum":                  6,
				"timesortfrom":              timesortfrom,
				"timesortto":                timesortto,
				"limittononsuspendedevents": true,
			},
		},
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(jsonBytes)
	resp, err := lms.PostData(fullReq, nil, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var dashResp []Dashboard
	err = json.NewDecoder(resp.Body).Decode(&dashResp)
	if err != nil {
		return nil, err
	}
	return &dashResp[0], nil
}
