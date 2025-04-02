package controller

import (
	"github.com/rainu/ask-mai/controller/history"
	"log/slog"
	"time"
)

func (c *Controller) saveHistory() {
	if c.appConfig.History.Path == "" {
		return
	}

	hw := history.NewWriter(c.appConfig.History.Path)
	we := hw.Write(historyMessagesToEntry(c.currentConversation))
	if we != nil {
		slog.Warn("Error writing history file!", "error", we.Error())
	}
}

func historyMessagesToEntry(messages LLMMessages) history.Entry {
	entry := history.Entry{
		Meta: history.EntryMeta{
			Version:   1,
			Timestamp: time.Now().UnixMilli(),
		},
		Content: history.EntryContent{
			Messages: make([]history.Message, len(messages)),
		},
	}

	for i, msg := range messages {
		entry.Content.Messages[i] = history.Message{
			Id:           msg.Id,
			Role:         msg.Role,
			ContentParts: make([]history.MessageContentPart, len(msg.ContentParts)),
		}
		for j, cp := range msg.ContentParts {
			entry.Content.Messages[i].ContentParts[j] = history.MessageContentPart{
				Type:    string(cp.Type),
				Content: cp.Content,
			}

			if cp.Call.Function != "" {
				entry.Content.Messages[i].ContentParts[j].Call = &history.MessageCall{
					Id:        cp.Call.Id,
					Function:  cp.Call.Function,
					Arguments: cp.Call.Arguments,
				}
				if cp.Call.Result != nil {
					entry.Content.Messages[i].ContentParts[j].Call.Result = &history.MessageCallResult{
						Content:    cp.Call.Result.Content,
						Error:      cp.Call.Result.Error,
						DurationMs: cp.Call.Result.DurationMs,
					}
				}
			}

		}
	}

	return entry
}

func historyEntry2Messages(entry history.Entry) LLMMessages {
	messages := make(LLMMessages, len(entry.Content.Messages))
	for i, msg := range entry.Content.Messages {
		messages[i] = LLMMessage{
			Id:           msg.Id,
			Role:         msg.Role,
			ContentParts: make([]LLMMessageContentPart, len(msg.ContentParts)),
		}
		for j, cp := range msg.ContentParts {
			messages[i].ContentParts[j] = LLMMessageContentPart{
				Type:    LLMMessageContentPartType(cp.Type),
				Content: cp.Content,
			}

			if cp.Call != nil {
				messages[i].ContentParts[j].Call = LLMMessageCall{
					Id:        cp.Call.Id,
					Function:  cp.Call.Function,
					Arguments: cp.Call.Arguments,
				}
				if cp.Call.Result != nil {
					messages[i].ContentParts[j].Call.Result = &LLMMessageCallResult{
						Content:    cp.Call.Result.Content,
						Error:      cp.Call.Result.Error,
						DurationMs: cp.Call.Result.DurationMs,
					}
				}
			}
		}
	}

	return messages
}

func (c *Controller) GetCurrentConversation() LLMMessages {
	return c.currentConversation
}
