import { mcp } from '../../../wailsjs/go/models.ts'

export type ToolCallResultContent = mcp.TextContent | mcp.ImageContent | mcp.AudioContent | mcp.EmbeddedResource

export type ToolCallResult = mcp.CallToolResult & {
	content: ToolCallResultContent[]
}