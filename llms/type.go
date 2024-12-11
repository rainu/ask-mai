package llms

import "github.com/tmc/langchaingo/llms"

type Model interface {
	llms.Model

	Close() error
}
