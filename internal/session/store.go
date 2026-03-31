package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a recorded wmux session.
type Entry struct {
	Name      string    `json:"name"`
	SessionID string    `json:"session_id"`
	CWD       string    `json:"cwd"`
	CreatedAt time.Time `json:"created_at"`
}

func storeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "wmux", "sessions")
	return dir, os.MkdirAll(dir, 0o755)
}

func entryPath(name string) (string, error) {
	dir, err := storeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".json"), nil
}

// Save records a session entry.
func Save(e Entry) error {
	path, err := entryPath(e.Name)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a single session entry by name.
func Load(name string) (Entry, error) {
	path, err := entryPath(name)
	if err != nil {
		return Entry{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, err
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, err
	}
	return e, nil
}

// List returns all recorded sessions, sorted by creation time (newest first).
func List() ([]Entry, error) {
	dir, err := storeDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var sessions []Entry
	for _, de := range entries {
		if filepath.Ext(de.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, de.Name()))
		if err != nil {
			continue
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			continue
		}
		sessions = append(sessions, e)
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})
	return sessions, nil
}

// Delete removes a session entry by name.
func Delete(name string) error {
	path, err := entryPath(name)
	if err != nil {
		return err
	}
	return os.Remove(path)
}
