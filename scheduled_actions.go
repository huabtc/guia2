package guia2

import (
	"encoding/json"
	"errors"
	"time"
)

// ScheduledActionStep describes a single step of a scheduled action.
type ScheduledActionStep struct {
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// ScheduledAction represents an Appium scheduled action request.
type ScheduledAction struct {
	Name            string
	Steps           []ScheduledActionStep
	Times           int
	Interval        time.Duration
	MaxHistoryItems int
	MaxPass         int
	MaxFail         int
}

// ActionStepException mirrors the exception payload returned by the server.
type ActionStepException struct {
	Name       string `json:"name"`
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace"`
}

// ActionStepResult contains an individual step execution result.
type ActionStepResult struct {
	Name      string               `json:"name"`
	Type      string               `json:"type"`
	Timestamp int64                `json:"timestamp"`
	Passed    bool                 `json:"passed"`
	Result    interface{}          `json:"result"`
	Exception *ActionStepException `json:"exception"`
}

// ActionHistory represents the response payload returned by action history endpoints.
type ActionHistory struct {
	Repeats     int64                `json:"repeats"`
	StepResults [][]ActionStepResult `json:"stepResults"`
}

// ScheduleAction registers a repeating action sequence on the device.
func (d *Driver) ScheduleAction(action ScheduledAction) error {
	if action.Name == "" {
		return errors.New("action name is required")
	}
	if len(action.Steps) == 0 {
		return errors.New("at least one action step is required")
	}
	payload := map[string]interface{}{
		"name":  action.Name,
		"steps": action.Steps,
	}
	if action.Times > 0 {
		payload["times"] = action.Times
	}
	if action.Interval > 0 {
		payload["intervalMs"] = action.Interval.Milliseconds()
	}
	if action.MaxHistoryItems > 0 {
		payload["maxHistoryItems"] = action.MaxHistoryItems
	}
	if action.MaxPass > 0 {
		payload["maxPass"] = action.MaxPass
	}
	if action.MaxFail > 0 {
		payload["maxFail"] = action.MaxFail
	}
	_, err := d.executePost(payload, "/session", d.sessionId, "appium/schedule_action")
	return err
}

// ActionHistory fetches the execution history for the given scheduled action.
func (d *Driver) ActionHistory(name string) (*ActionHistory, error) {
	return d.fetchActionHistory("appium/action_history", name)
}

// UnscheduleAction removes a scheduled action and returns the final execution history.
func (d *Driver) UnscheduleAction(name string) (*ActionHistory, error) {
	return d.fetchActionHistory("appium/unschedule_action", name)
}

func (d *Driver) fetchActionHistory(pathSuffix string, name string) (*ActionHistory, error) {
	if name == "" {
		return nil, errors.New("action name is required")
	}
	payload := map[string]interface{}{
		"name": name,
	}
	rawResp, err := d.executePost(payload, "/session", d.sessionId, pathSuffix)
	if err != nil {
		return nil, err
	}
	var reply struct {
		Value ActionHistory `json:"value"`
	}
	if err := json.Unmarshal(rawResp, &reply); err != nil {
		return nil, err
	}
	return &reply.Value, nil
}
