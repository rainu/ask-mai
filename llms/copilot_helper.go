package llms

import cmdchain "github.com/rainu/go-command-chain"

var copilotInstalled *bool

func IsCopilotInstalled() bool {
	if copilotInstalled == nil {
		err := cmdchain.Builder().
			Join("gh", "copilot", "-v").
			Finalize().Run()

		installed := err == nil
		copilotInstalled = &installed
	}

	return *copilotInstalled
}
