package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SettingsPath returns the path to ~/.claude/settings.json.
func SettingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "settings.json")
}

// hookEntry matches the structure in settings.json hooks.
type hookMatcher struct {
	Matcher string       `json:"matcher,omitempty"`
	Hooks   []hookAction `json:"hooks"`
}

type hookAction struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// wzmuxCommand returns the hook command string for a given event.
func wzmuxCommand(event string) string {
	return fmt.Sprintf("wzmux hook %s", event)
}

// isWmuxHook checks if a hook action belongs to wzmux.
func isWmuxHook(action hookAction) bool {
	return strings.Contains(action.Command, "wzmux hook")
}

// hookEvents are the Claude Code events wzmux registers for.
var hookEvents = []string{"PreToolUse", "PostToolUse", "Stop", "Notification"}

// AddHooks adds wzmux hook entries to settings.json, preserving existing hooks.
func AddHooks() error {
	path := SettingsPath()
	settings, err := readSettings(path)
	if err != nil {
		return err
	}

	// Ensure hooks map exists
	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		hooks = map[string]interface{}{}
		settings["hooks"] = hooks
	}

	for _, event := range hookEvents {
		addEventHook(hooks, event)
	}

	return writeSettings(path, settings)
}

// RemoveHooks removes wzmux hook entries from settings.json.
func RemoveHooks() error {
	path := SettingsPath()
	settings, err := readSettings(path)
	if err != nil {
		return err
	}

	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		return nil // no hooks to remove
	}

	for _, event := range hookEvents {
		removeEventHook(hooks, event)
	}

	// Remove hooks key if empty
	if len(hooks) == 0 {
		delete(settings, "hooks")
	}

	return writeSettings(path, settings)
}

func addEventHook(hooks map[string]interface{}, event string) {
	wzmuxCmd := wzmuxCommand(event)

	// Get existing matchers for this event
	var matchers []interface{}
	if existing, ok := hooks[event]; ok {
		if arr, ok := existing.([]interface{}); ok {
			matchers = arr
		}
	}

	// Check if wzmux hook already exists
	for _, m := range matchers {
		mMap, ok := m.(map[string]interface{})
		if !ok {
			continue
		}
		hookActions, ok := mMap["hooks"].([]interface{})
		if !ok {
			continue
		}
		for _, h := range hookActions {
			hMap, ok := h.(map[string]interface{})
			if !ok {
				continue
			}
			if cmd, ok := hMap["command"].(string); ok && strings.Contains(cmd, "wzmux hook") {
				return // already registered
			}
		}
	}

	// Add wzmux hook entry
	newMatcher := map[string]interface{}{
		"hooks": []interface{}{
			map[string]interface{}{
				"type":    "command",
				"command": wzmuxCmd,
			},
		},
	}
	matchers = append(matchers, newMatcher)
	hooks[event] = matchers
}

func removeEventHook(hooks map[string]interface{}, event string) {
	matchers, ok := hooks[event].([]interface{})
	if !ok {
		return
	}

	var kept []interface{}
	for _, m := range matchers {
		mMap, ok := m.(map[string]interface{})
		if !ok {
			kept = append(kept, m)
			continue
		}
		hookActions, ok := mMap["hooks"].([]interface{})
		if !ok {
			kept = append(kept, m)
			continue
		}
		var keptActions []interface{}
		for _, h := range hookActions {
			hMap, ok := h.(map[string]interface{})
			if !ok {
				keptActions = append(keptActions, h)
				continue
			}
			if cmd, ok := hMap["command"].(string); ok && strings.Contains(cmd, "wzmux hook") {
				continue // remove this one
			}
			keptActions = append(keptActions, h)
		}
		if len(keptActions) > 0 {
			mMap["hooks"] = keptActions
			kept = append(kept, mMap)
		}
	}

	if len(kept) > 0 {
		hooks[event] = kept
	} else {
		delete(hooks, event)
	}
}

func readSettings(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("read settings: %w", err)
	}
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("parse settings: %w", err)
	}
	return settings, nil
}

func writeSettings(path string, settings map[string]interface{}) error {
	// Backup before writing
	if _, err := os.Stat(path); err == nil {
		data, _ := os.ReadFile(path)
		os.WriteFile(path+".bak", data, 0644)
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
