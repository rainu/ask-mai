package approval

import (
	"context"
	"encoding/json"
	"github.com/rainu/ask-mai/internal/expression"
	"log/slog"
	"strings"
)

type Approval string

const (
	Always = "always"
	Never  = "never"
)

type Variables struct {
	ToolDefinition  any    `json:"definition"`
	RawArguments    string `json:"raw_args"`
	ParsedArguments any    `json:"args"`
}

func (a Approval) NeedsApproval(ctx context.Context, jsonArgs string, td any) bool {
	if a == "" {
		// No approval expression is set, so we assume no approval is needed
		return false
	}
	switch strings.TrimSpace(strings.ToLower(string(a))) {
	case Always:
		return true
	case Never:
		return false
	}

	exVars := Variables{
		RawArguments: jsonArgs,
	}
	if td != nil {
		exVars.ToolDefinition = td
	}

	err := json.Unmarshal([]byte(jsonArgs), &exVars.ParsedArguments)
	if err != nil {
		slog.Warn("error parsing arguments", "args", jsonArgs, "error", err)
	}

	b, err := expression.Run(ctx, string(a), exVars).AsBoolean()

	if err != nil {
		slog.Error("error running approval expression", "expression", string(a), "error", err)
		return true
	}
	return b
}
