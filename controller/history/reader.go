package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kteru/reversereader"
	"io"
	"os"
	"slices"
	"strings"
)

type Reader struct {
	historyPath string
}

func NewReader(historyPath string) *Reader {
	return &Reader{
		historyPath: historyPath,
	}
}

func (r *Reader) GetCount() (int, error) {
	count := 0
	err := r.readLines(0, -1, func(line string) error {
		count++
		return nil
	})

	return count, err
}

func (r *Reader) GetLast(skip, limit int) ([]Entry, error) {
	entries := make([]Entry, 0, limit)
	err := r.readLines(skip, limit, func(line string) error {
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return fmt.Errorf("failed to unmarshal history entry: %s", err)
		}
		entries = append(entries, entry)

		return nil
	})

	return entries, err
}

func (r *Reader) Search(matcher func(entry Entry) (bool, bool)) ([]Entry, error) {
	var entries []Entry
	err := r.readLines(0, -1, func(line string) error {
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return fmt.Errorf("failed to unmarshal history entry: %s", err)
		}

		matches, continue_ := matcher(entry)
		if matches {
			entries = append(entries, entry)
		}
		if !continue_ {
			return io.EOF
		}

		return nil
	})

	return entries, err
}

func (r *Reader) readLines(skip, limit int, processor func(line string) error) error {
	if r.historyPath == "" {
		return nil
	}

	file, err := os.Open(r.historyPath)
	if err != nil {
		return fmt.Errorf("failed to open history file %s: %s", r.historyPath, err)
	}
	defer file.Close()

	rr := reversereader.NewReader(file)

	scanner := bufio.NewScanner(rr)
	lineCount := 0
	processed := 0
	skipped := 0

	for scanner.Scan() {
		line := scanner.Bytes()

		// because we are reading the file in reverse, we need to reverse the line
		slices.Reverse(line)

		sLine := string(line)
		if strings.TrimSpace(sLine) == "" {
			continue // Skip empty lines
		}

		if skip > skipped {
			skipped++
			continue
		}

		if err := processor(sLine); err != nil {
			if err == io.EOF {
				break // Stop processing if EOF is reached
			}
			return fmt.Errorf("error processing line %d: %w", lineCount, err)
		}
		processed++

		if limit > 0 && processed >= limit {
			break
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading history file: %s", err)
	}

	return nil
}
