package validate

import (
	"fmt"
	"strings"

	"boxed/internal/box"
)

// BoxType validates that a box type string is one of the supported types.
// This is the first line of defense for fail-fast error handling, rejecting
// invalid input before any rendering work begins.
func BoxType(t string) error {
	boxType := box.BoxType(t)
	if !boxType.IsValid() {
		validTypes := make([]string, 0, 4)
		for _, bt := range box.AllBoxTypes() {
			validTypes = append(validTypes, bt.String())
		}
		return fmt.Errorf("invalid box type %q, must be one of: %s",
			t, strings.Join(validTypes, ", "))
	}
	return nil
}

// KVPair validates that a key-value string is properly formatted as "key=value".
// Enforces that keys cannot be empty, as empty keys would be meaningless in the
// rendered output. Values can be empty to support cases like "flag=" indicating
// an unset or empty value.
func KVPair(kv string) error {
	if !strings.Contains(kv, "=") {
		return fmt.Errorf("invalid key-value pair %q: must be in format key=value", kv)
	}

	parts := strings.SplitN(kv, "=", 2)
	if parts[0] == "" {
		return fmt.Errorf("invalid key-value pair %q: key cannot be empty", kv)
	}

	return nil
}

// BorderStyle validates that a border style string is supported by Lip Gloss.
// While Lip Gloss would gracefully handle unknown styles by falling back to a default,
// we validate explicitly to provide clear error messages to users at parse time.
func BorderStyle(style string) error {
	validStyles := map[string]bool{
		"":       true,
		"normal": true,
		"rounded": true,
		"thick":  true,
		"double": true,
	}

	if !validStyles[style] {
		return fmt.Errorf("invalid border style %q, must be one of: normal, rounded, thick, double", style)
	}
	return nil
}

// Box performs final validation on a complete box model before rendering.
// This catches logical errors that pass individual field validation but create
// invalid combinations, like a box with no displayable content.
func Box(b *box.Box) error {
	if !b.HasContent() {
		return fmt.Errorf("box has no content: provide at least one of --title, --subtitle, --kv, or --footer")
	}

	if b.Width < 0 {
		return fmt.Errorf("width must be non-negative, got %d", b.Width)
	}

	return nil
}
