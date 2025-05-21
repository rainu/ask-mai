package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rainu/ask-mai/internal/llms/tools"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type FileReadingArguments struct {
	Path   Path               `json:"path"`
	Limits *FileReadingLimits `json:"limits"`
}

type FileReadingLimits struct {
	Mode   string `json:"m"`
	Offset int    `json:"o"`
	Limit  int    `json:"l"`
}

type FileReadingResult struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

var FileReadingDefinition = tools.BuiltinDefinition{
	Description: "Reading a text file from the user's system.",
	Parameter: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "The absolute path to the text file to be read. Use '~' as placeholder for the user's home directory.",
			},
			"limits": map[string]any{
				"type":        "object",
				"description": "The (optional) limits for reading the file.",
				"properties": map[string]any{
					"m": map[string]any{
						"type":        "string",
						"description": "The limit mode.",
						"enum":        []string{"line", "character"},
					},
					"o": map[string]any{
						"type":        "number",
						"description": "The line or char offset to start reading from. Default is 0.",
					},
					"l": map[string]any{
						"type":        "number",
						"description": "The line or char limit to read. Default is -1 (read all).",
					},
				},
				"additionalProperties": false,
				"required":             []string{"m"},
			},
		},
		"additionalProperties": false,
		"required":             []string{"path"},
	},
	Function: func(ctx context.Context, jsonArguments string) ([]byte, error) {
		var pArgs FileReadingArguments
		err := json.Unmarshal([]byte(jsonArguments), &pArgs)
		if err != nil {
			return nil, fmt.Errorf("error parsing arguments: %w", err)
		}

		if string(pArgs.Path) == "" {
			return nil, fmt.Errorf("missing parameter: 'path'")
		}
		path, err := pArgs.Path.Get()
		if err != nil {
			return nil, err
		}

		mode := ""
		if pArgs.Limits != nil {
			if pArgs.Limits.Mode != "line" && pArgs.Limits.Mode != "char" {
				return nil, fmt.Errorf("invalid limit mode: '%s'", pArgs.Limits.Mode)
			}
			if pArgs.Limits.Limit <= -1 {
				pArgs.Limits.Limit = -1
			}
			if pArgs.Limits.Offset < 0 {
				pArgs.Limits.Offset = 0
			}
			mode = pArgs.Limits.Mode
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}
		defer file.Close()

		var content string
		switch mode {
		case "line":
			// read line by line
			content, err = readLines(file, pArgs.Limits.Offset, pArgs.Limits.Limit)
		case "char":
			// read char by char
			content, err = readChars(file, pArgs.Limits.Offset, pArgs.Limits.Limit)
		default:
			content, err = readAll(file)
		}

		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		absolutePath, err := filepath.Abs(file.Name())
		if err != nil {
			slog.Warn("Error getting absolute path!", "error", err)
			absolutePath = file.Name()
		}

		return json.Marshal(FileReadingResult{
			Path:    absolutePath,
			Content: content,
		})
	},
}

func readLines(file *os.File, offset, limit int) (string, error) {
	content := strings.Builder{}
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		if i < offset {
			continue
		}
		if limit != -1 && i >= offset+limit {
			break
		}
		content.WriteString(scanner.Text() + "\n")
	}
	return content.String(), scanner.Err()
}

func readChars(file *os.File, offset int, limit int) (string, error) {
	content := strings.Builder{}
	reader := bufio.NewReader(file)
	for i := 0; i < offset; i++ {
		if _, _, err := reader.ReadRune(); err != nil {
			return content.String(), err
		}
	}
	for i := 0; limit == -1 || i < limit; i++ {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return content.String(), err
		}
		content.WriteRune(r)
	}
	return content.String(), nil
}

func readAll(file *os.File) (string, error) {
	c, e := io.ReadAll(file)
	return string(c), e
}
