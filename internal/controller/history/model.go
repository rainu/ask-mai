package history

import "time"

type Entry struct {
	Meta    EntryMeta    `json:"m"`
	Content EntryContent `json:"c"`
}

type EntryMeta struct {
	Version   int   `json:"v"`
	Timestamp int64 `json:"t"`
}

type EntryContent struct {
	Messages []Message `json:"m"`
}

type Message struct {
	Id           string               `json:"i,omitempty"`
	Role         string               `json:"r,omitempty"`
	ContentParts []MessageContentPart `json:"p,omitempty"`
	CreatedAt    int64                `json:"t,omitempty"`
}
type MessageContentPart struct {
	Type    string       `json:"t,omitempty"`
	Content string       `json:"c,omitempty"`
	Call    *MessageCall `json:"ca,omitempty"`
}

type MessageCall struct {
	Id        string             `json:"i,omitempty"`
	Function  string             `json:"f,omitempty"`
	Arguments string             `json:"a,omitempty"`
	Result    *MessageCallResult `json:"r,omitempty"`
}

type MessageCallResult struct {
	Content    string `json:"c"`
	Error      string `json:"e"`
	DurationMs int64  `json:"d"`
}

func NewEntry(content EntryContent) Entry {
	return Entry{
		Meta: EntryMeta{
			Version:   1,
			Timestamp: time.Now().UnixMilli(),
		},
		Content: content,
	}
}
