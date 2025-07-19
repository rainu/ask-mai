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
	CommandLine           string            `json:"command,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Arguments             []string          `json:"arguments,omitempty"`
	Environment           map[string]string `json:"env"`
	AdditionalEnvironment map[string]string `json:"additionalEnv"`
	WorkingDirectory      string            `json:"workingDir"`
	Output                *OutputSettings   `json:"output,omitempty"`
}

type OutputSettings struct {
	DisableStdOut bool `json:"disableStdOut"`
	DisableStdErr bool `json:"disableStdErr"`
	FirstNBytes   int  `json:"firstNBytes"`
	LastNBytes    int  `json:"lastNBytes"`
}

func (c CommandDescriptor) Run(ctx context.Context) ([]byte, error) {
	var cmdBuild cmdchain.CommandBuilder

	if c.CommandLine != "" {
		cmdBuild = cmdchain.Builder().JoinShellCmdWithContext(ctx, c.CommandLine)
	} else {
		cmdBuild = cmdchain.Builder().JoinWithContext(ctx, c.Name, c.Arguments...)
	}

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
	if c.Output == nil || !c.Output.DisableStdOut {
		cmd = cmd.WithOutput(oFile)
	}
	if c.Output == nil || !c.Output.DisableStdErr {
		cmd = cmd.WithError(oFile)
	}

	execErr := cmd.Run()
	return c.getOutput(oFile), execErr
}

func (c CommandDescriptor) getOutput(f *os.File) []byte {
	if c.Output == nil {
		return readFile(f)
	}
	cfg := c.Output
	if cfg.FirstNBytes < 0 || cfg.LastNBytes < 0 {
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
	if cfg.FirstNBytes+cfg.LastNBytes > int(fs.Size()) {
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

	if cfg.FirstNBytes > 0 {
		_, err = io.CopyN(buf, f, int64(cfg.FirstNBytes))
		if err != nil && err != io.EOF {
			slog.Error("Could not read first bytes from command output file.",
				"bytes", cfg.FirstNBytes,
				"path", f.Name(),
				"error", err,
			)
			return nil
		}
	} else {
		// Indicate that there were bytes skipped
		buf.WriteString(skippedBytesIndicator(fs.Size() - int64(cfg.LastNBytes)))
		buf.WriteRune('\n')
	}

	if cfg.FirstNBytes > 0 && cfg.LastNBytes > 0 {
		// Indicate that there were bytes skipped
		buf.WriteRune('\n')
		buf.WriteString(skippedBytesIndicator(fs.Size() - int64(cfg.FirstNBytes+cfg.LastNBytes)))
		buf.WriteRune('\n')
	}

	if cfg.LastNBytes > 0 {
		_, err = f.Seek(-int64(cfg.LastNBytes), io.SeekEnd) // Seek to the last N bytes
		if err != nil {
			slog.Error("Could not seek to the last bytes of command output file.",
				"bytes", cfg.LastNBytes,
				"path", f.Name(),
				"error", err,
			)
			return nil
		}
		_, err = io.Copy(buf, f)
		if err != nil && err != io.EOF {
			slog.Error("Could not read last bytes from command output file.",
				"bytes", cfg.LastNBytes,
				"path", f.Name(),
				"error", err,
			)
			return nil
		}
	} else {
		// Indicate that there were bytes skipped
		buf.WriteRune('\n')
		buf.WriteString(skippedBytesIndicator(fs.Size() - int64(cfg.FirstNBytes)))
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
