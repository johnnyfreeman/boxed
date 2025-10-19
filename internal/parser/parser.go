package parser

import (
	"strings"

	"boxed/internal/box"
	"boxed/internal/validate"
)

// Options represents the parsed CLI flags for box configuration.
// This struct acts as an intermediate representation between Cobra's
// flag parsing and our Box model, keeping CLI concerns separate from
// domain logic.
type Options struct {
	Title       string
	Subtitle    string
	KVFlags     []string
	Footer      string
	Width       int
	BorderStyle string
}

// ParseBox converts CLI arguments into a validated Box model.
// This is a pure function that performs all validation upfront (fail-fast)
// before constructing the box, ensuring that any Box instance that successfully
// returns from this function is guaranteed to be valid and renderable.
func ParseBox(boxType string, opts Options) (*box.Box, error) {
	if err := validate.BoxType(boxType); err != nil {
		return nil, err
	}

	if err := validate.BorderStyle(opts.BorderStyle); err != nil {
		return nil, err
	}

	kvPairs, err := parseKVPairs(opts.KVFlags)
	if err != nil {
		return nil, err
	}

	b := &box.Box{
		Type:        box.BoxType(boxType),
		Title:       opts.Title,
		Subtitle:    opts.Subtitle,
		KVPairs:     kvPairs,
		Footer:      opts.Footer,
		Width:       opts.Width,
		BorderStyle: opts.BorderStyle,
	}

	if err := validate.Box(b); err != nil {
		return nil, err
	}

	return b, nil
}

// parseKVPairs converts an array of "key=value" strings into KV structs.
// Each string is validated before parsing to ensure fail-fast behavior.
// Supports comma-separated pairs (e.g., "A=1,B=2,C=3") for convenience.
// Splits on the first '=' only, allowing '=' characters in values
// (e.g., "url=http://example.com?a=1&b=2").
func parseKVPairs(kvFlags []string) ([]box.KV, error) {
	if len(kvFlags) == 0 {
		return nil, nil
	}

	kvPairs := make([]box.KV, 0, len(kvFlags))
	for _, kvFlag := range kvFlags {
		// Split on comma to support comma-separated pairs like "A=1,B=2,C=3"
		parts := strings.Split(kvFlag, ",")

		for _, kv := range parts {
			kv = strings.TrimSpace(kv)

			if kv == "" {
				continue
			}

			if err := validate.KVPair(kv); err != nil {
				return nil, err
			}

			pair := strings.SplitN(kv, "=", 2)
			kvPairs = append(kvPairs, box.KV{
				Key:   pair[0],
				Value: pair[1],
			})
		}
	}

	return kvPairs, nil
}
