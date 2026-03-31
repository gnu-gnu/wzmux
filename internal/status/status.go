package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const StatusDir = "/tmp/claude-agent-status"

// AgentStatus represents the JSON status file for an agent.
type AgentStatus struct {
	Status string `json:"status"`
	Pane   string `json:"pane"`
	CWD    string `json:"cwd"`
	TS     int64  `json:"ts"`
	Msg    string `json:"msg"`
}

// Read reads the status file for a given pane ID.
func Read(paneID int) (*AgentStatus, error) {
	path := filepath.Join(StatusDir, fmt.Sprintf("pane-%d.json", paneID))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s AgentStatus
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Write writes a status file for the given pane ID.
func Write(paneID int, statusStr, cwd, msg string) error {
	if err := os.MkdirAll(StatusDir, 0755); err != nil {
		return err
	}
	s := AgentStatus{
		Status: statusStr,
		Pane:   fmt.Sprintf("%d", paneID),
		CWD:    cwd,
		TS:     time.Now().Unix(),
		Msg:    msg,
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	path := filepath.Join(StatusDir, fmt.Sprintf("pane-%d.json", paneID))
	return os.WriteFile(path, data, 0644)
}

// Remove removes the status file for a given pane ID.
func Remove(paneID int) error {
	path := filepath.Join(StatusDir, fmt.Sprintf("pane-%d.json", paneID))
	return os.Remove(path)
}

// CleanAll removes all status files.
func CleanAll() error {
	return os.RemoveAll(StatusDir)
}
