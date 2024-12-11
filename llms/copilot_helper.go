package llms

import cmdchain "github.com/rainu/go-command-chain"

func IsCopilotInstalled() bool {
	err := cmdchain.Builder().
		Join("gh", "copilot", "-v").
		Finalize().Run()

	return err == nil
}
