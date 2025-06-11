//go:build linux

package system

import (
	cmdchain "github.com/rainu/go-command-chain"
	"sync"
)

type OSInfo struct {
	LSBRelease string `json:"lsb_release,omitempty"`
	Release    string `json:"release,omitempty"`
	Uname      string `json:"uname,omitempty"`
}

func init() {
	getOSInfo = func() any {
		result := OSInfo{}
		wg := sync.WaitGroup{}

		result.LSBRelease = exec(&wg, "lsb_release", "-a")
		result.Release = exec(&wg, "sh", "-c", "cat /etc/*release")
		result.Uname = exec(&wg, "uname", "-a")

		wg.Wait()
		return result
	}
}

func exec(wg *sync.WaitGroup, cmd string, args ...string) string {
	wg.Add(1)

	result := make(chan string)
	go func() {
		defer wg.Done()
		defer close(result)

		r, _, _ := cmdchain.Builder().Join(cmd, args...).Finalize().RunAndGet()
		result <- r
	}()

	return <-result
}
