package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type CallArguments struct {
	Method string            `json:"method"`
	Url    string            `json:"url"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

var DefaultClient *http.Client

func init() {
	DefaultClient = &http.Client{}
	DefaultClient.Jar, _ = cookiejar.New(nil)
}

var CallTool = mcp.NewTool("callHttp",
	mcp.WithDescription("Do a http call to a given url with a given method and body."),
	mcp.WithString("method",
		mcp.Enum(
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		),
		mcp.Description("The method to use. Default is "+http.MethodGet+"."),
	),
	mcp.WithString("url",
		mcp.Required(),
		mcp.Description("The url to call."),
	),
	mcp.WithObject("header",
		mcp.Description("The headers to send."),
		mcp.AdditionalProperties(map[string]any{"additionalProperties": true}),
	),
	mcp.WithString("body",
		mcp.Description("The body to send. Default is empty."),
	),
)

var CallToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var pArgs CallArguments

	r, w := io.Pipe()
	go func() {
		defer w.Close()

		json.NewEncoder(w).Encode(request.Params.Arguments)
	}()

	err := json.NewDecoder(r).Decode(&pArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	callDesc := CallDescriptor{
		Method: pArgs.Method,
		Url:    pArgs.Url,
		Header: pArgs.Header,
		Body:   strings.NewReader(pArgs.Body),
	}

	result, err := callDesc.Run(ctx, DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("error executing call: %w", err)
	}

	raw, err := json.Marshal(result)
	return mcp.NewToolResultText(string(raw)), err
}
