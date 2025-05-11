package controller

import (
	"github.com/rainu/ask-mai/controller/history"
	"log/slog"
	"regexp"
	"strings"
)

func (c *Controller) saveHistory() {
	if *c.getProfile().History.Path == "" {
		return
	}
	if len(c.currentConversation) == 0 {
		// prevent writing empty entries
		return
	}

	hw := history.NewWriter(*c.getProfile().History.Path)
	we := hw.Write(historyMessagesToEntry(c.currentConversation))
	if we != nil {
		slog.Warn("Error writing history file!", "error", we.Error())
	}
}

func historyMessagesToEntry(messages LLMMessages) history.Entry {
	entry := history.NewEntry(history.EntryContent{
		Messages: make([]history.Message, len(messages)),
	})

	for i, msg := range messages {
		entry.Content.Messages[i] = history.Message{
			Id:           msg.Id,
			Role:         msg.Role,
			ContentParts: make([]history.MessageContentPart, len(msg.ContentParts)),
			CreatedAt:    msg.Created,
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
			Created:      msg.CreatedAt,
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

func (c *Controller) HistoryGetCount() (int, error) {
	if *c.getProfile().History.Path == "" {
		return 0, nil
	}
	hr := history.NewReader(*c.getProfile().History.Path)

	return hr.GetCount()
}

func (c *Controller) HistoryGetLast(skip, limit int) ([]history.Entry, error) {
	if *c.getProfile().History.Path == "" {
		return nil, nil
	}

	hr := history.NewReader(*c.getProfile().History.Path)

	return hr.GetLast(skip, limit)
}

func (c *Controller) HistorySearch(query string) ([]history.Entry, error) {
	if *c.getProfile().History.Path == "" {
		return nil, nil
	}

	queryExp, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	hr := history.NewReader(*c.getProfile().History.Path)

	return hr.Search(func(entry history.Entry) (bool, bool) {
		content := strings.Builder{}
		for _, message := range entry.Content.Messages {
			for _, part := range message.ContentParts {
				content.WriteString(part.Content)
				if part.Call != nil {
					content.WriteString(part.Call.Function)
					content.WriteString("(")
					content.WriteString(part.Call.Arguments)
					content.WriteString(")")
					if part.Call.Result != nil {
						content.WriteString(" -> ")
						content.WriteString(part.Call.Result.Content)
					}
				}
			}
		}

		return queryExp.MatchString(content.String()), true
	})
}
