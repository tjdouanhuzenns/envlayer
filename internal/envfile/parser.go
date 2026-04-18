package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed environment file as an ordered map.
type EnvMap struct {
	Keys   []string
	Values map[string]string
}

func NewEnvMap() *EnvMap {
	return &EnvMap{Values: make(map[string]string)}
}

// Set adds or updates a key.
func (e *EnvMap) Set(key, value string) {
	if _, exists := e.Values[key]; !exists {
		e.Keys = append(e.Keys, key)
	}
	e.Values[key] = value
}

// ParseFile reads and parses a .env file into an EnvMap.
func ParseFile(path string) (*EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	env := NewEnvMap()
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s:%d: invalid line: %q", path, lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		env.Set(key, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning %s: %w", path, err)
	}

	return env, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
