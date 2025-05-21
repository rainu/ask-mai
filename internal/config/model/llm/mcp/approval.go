package mcp

import (
	"context"
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rainu/ask-mai/internal/expression"
	"log/slog"
	"strings"
)

type Approval string

const (
	ApprovalAlways = "true"
	ApprovalNever  = "false"
)

type ApprovalVariables struct {
	ToolDefinition  mcp.Tool `json:"definition"`
	RawArguments    string   `json:"raw_args"`
	ParsedArguments any      `json:"args"`
}

func (a Approval) NeedsApproval(ctx context.Context, jsonArgs string, td *mcp.Tool) bool {
	if a == "" {
		// No approval expression is set, so we assume no approval is needed
		return false
	}
	switch strings.TrimSpace(strings.ToLower(string(a))) {
	case ApprovalAlways:
		return true
	case ApprovalNever:
		return false
	}

	exVars := ApprovalVariables{
		RawArguments: jsonArgs,
	}
	if td != nil {
		exVars.ToolDefinition = *td
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
