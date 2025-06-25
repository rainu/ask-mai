package command

import (
	"bytes"
	"context"
	"fmt"
	cmdchain "github.com/rainu/go-command-chain"
	"io"
	"log/slog"
	"os"
)

type CommandDescriptor struct {
	Command               string            `json:"command"`
	Arguments             []string          `json:"arguments"`
	Environment           map[string]string `json:"env"`
	AdditionalEnvironment map[string]string `json:"additionalEnv"`
	WorkingDirectory      string            `json:"workingDir"`
	DisableStdOut         bool              `json:"disableStdOut"`
	DisableStdErr         bool              `json:"disableStdErr"`
	FirstNBytes           int               `json:"firstNBytes"`
	LastNBytes            int               `json:"lastNBytes"`
}

func (c CommandDescriptor) Run(ctx context.Context) ([]byte, error) {
	cmdBuild := cmdchain.Builder().JoinWithContext(ctx, c.Command, c.Arguments...)

	if len(c.Environment) > 0 {
		cmdBuild = cmdBuild.WithEnvironmentMap(toAnyMap(c.Environment))
	}
	if len(c.AdditionalEnvironment) > 0 {
		cmdBuild = cmdBuild.WithAdditionalEnvironmentMap(toAnyMap(c.AdditionalEnvironment))
	}
	if c.WorkingDirectory != "" {
		cmdBuild = cmdBuild.WithWorkingDirectory(c.WorkingDirectory)
	}

	oFile, err := os.CreateTemp("", "ask-mai.mcp.command.*")
	if err != nil {
		return nil, fmt.Errorf("could not create temporary file: %w", err)
	}
	defer func() {
		oFile.Close()
		os.Remove(oFile.Name())
	}()

	cmd := cmdBuild.WithErrorChecker(cmdchain.IgnoreExitErrors()).Finalize()
	if !c.DisableStdOut {
		cmd = cmd.WithOutput(oFile)
	}
	if !c.DisableStdErr {
		cmd = cmd.WithError(oFile)
	}

	execErr := cmd.Run()
	return c.getOutput(oFile), execErr
}

func (c CommandDescriptor) getOutput(f *os.File) []byte {
	if c.FirstNBytes < 0 || c.LastNBytes < 0 {
		return readFile(f)
	}

	fs, err := f.Stat()
	if err != nil {
		slog.Error("Could not get stats from command output file.",
			"path", f.Name(),
			"error", err,
		)
		return nil
	}
	if c.FirstNBytes+c.LastNBytes > int(fs.Size()) {
		return readFile(f)
	}

	buf := bytes.NewBuffer(nil)
	_, err = f.Seek(0, 0) // Reset file pointer to the beginning
	if err != nil {
		slog.Error("Could not seek to the beginning of command output file.",
			"path", f.Name(),
			"error", err,
		)
		return nil
	}

	if c.FirstNBytes > 0 {
		_, err = io.CopyN(buf, f, int64(c.FirstNBytes))
		if err != nil && err != io.EOF {
			slog.Error("Could not read first bytes from command output file.",
				"bytes", c.FirstNBytes,
				"path", f.Name(),
				"error", err,
			)
			return nil
		}
	} else {
		// Indicate that there were bytes skipped
		buf.WriteString(skippedBytesIndicator(fs.Size() - int64(c.LastNBytes)))
		buf.WriteRune('\n')
	}

	if c.FirstNBytes > 0 && c.LastNBytes > 0 {
		// Indicate that there were bytes skipped
		buf.WriteRune('\n')
		buf.WriteString(skippedBytesIndicator(fs.Size() - int64(c.FirstNBytes+c.LastNBytes)))
		buf.WriteRune('\n')
	}

	if c.LastNBytes > 0 {
		_, err = f.Seek(-int64(c.LastNBytes), io.SeekEnd) // Seek to the last N bytes
		if err != nil {
			slog.Error("Could not seek to the last bytes of command output file.",
				"bytes", c.LastNBytes,
				"path", f.Name(),
				"error", err,
			)
			return nil
		}
		_, err = io.Copy(buf, f)
		if err != nil && err != io.EOF {
			slog.Error("Could not read last bytes from command output file.",
				"bytes", c.LastNBytes,
				"path", f.Name(),
				"error", err,
			)
			return nil
		}
	} else {
		// Indicate that there were bytes skipped
		buf.WriteRune('\n')
		buf.WriteString(skippedBytesIndicator(fs.Size() - int64(c.FirstNBytes)))
	}

	return buf.Bytes()
}

func skippedBytesIndicator(skipped int64) string {
	return fmt.Sprintf("{{ %d bytes skipped }}", skipped)
}

func readFile(f *os.File) []byte {
	content, err := os.ReadFile(f.Name())
	if err != nil {
		slog.Error("Could not read file.",
			"path", f.Name(),
			"error", err,
		)
		return nil
	}
	return content
}

func toAnyMap(m map[string]string) map[any]any {
	result := map[any]any{}
	for k, v := range m {
		result[k] = v
	}
	return result
}
