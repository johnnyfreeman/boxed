package io

import (
	"bufio"
	"io"
	"strings"

	"boxed/internal/box"
	"boxed/internal/validate"
)

// KVReader abstracts reading key-value pairs from an input source.
// This interface enables dependency injection, allowing tests to use in-memory
// readers while production code reads from os.Stdin.
type KVReader interface {
	ReadKVPairs() ([]box.KV, error)
}

// StdinKVReader reads key-value pairs from an io.Reader (typically stdin).
// Each line should be in "key=value" format. The reader performs validation
// on each line as it's read, failing fast on the first invalid line rather
// than continuing to parse potentially malformed input.
type StdinKVReader struct {
	reader io.Reader
}

// NewStdinKVReader creates a reader that parses KV pairs from the given input.
// The io.Reader interface allows testing with strings.Reader while production
// uses os.Stdin, following the "accept interfaces" principle for maximum testability.
func NewStdinKVReader(r io.Reader) *StdinKVReader {
	return &StdinKVReader{reader: r}
}

// ReadKVPairs reads all lines from the input, parsing each as a key-value pair.
// Empty lines and lines containing only whitespace are skipped to handle
// trailing newlines and user formatting gracefully. The function validates each
// line before parsing to ensure fail-fast behavior.
func (s *StdinKVReader) ReadKVPairs() ([]box.KV, error) {
	var kvPairs []box.KV

	scanner := bufio.NewScanner(s.reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if err := validate.KVPair(line); err != nil {
			return nil, err
		}

		parts := strings.SplitN(line, "=", 2)
		kvPairs = append(kvPairs, box.KV{
			Key:   parts[0],
			Value: parts[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return kvPairs, nil
}
