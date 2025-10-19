package box

import "fmt"

// BoxType defines semantic meaning for terminal output boxes, driving both
// visual styling (color) and user interpretation. Using a constrained type
// rather than free-form strings ensures fail-fast validation at parse time
// and prevents typos from reaching the renderer.
type BoxType string

const (
	Success BoxType = "success"
	Error   BoxType = "error"
	Info    BoxType = "info"
	Warning BoxType = "warning"
)

func (b BoxType) String() string {
	return string(b)
}

func (b BoxType) IsValid() bool {
	switch b {
	case Success, Error, Info, Warning:
		return true
	default:
		return false
	}
}

func AllBoxTypes() []BoxType {
	return []BoxType{Success, Error, Info, Warning}
}

// KV represents key-value metadata displayed in the box content area.
// The String() method exists primarily for testing and debugging; the actual
// rendering logic uses Key and Value fields directly to allow flexible formatting.
type KV struct {
	Key   string
	Value string
}

func (kv KV) String() string {
	return fmt.Sprintf("%s=%s", kv.Key, kv.Value)
}

// Box is the core data model representing all content and configuration for
// a single terminal box render. It intentionally separates data (what to display)
// from presentation (how to display it), enabling dependency injection of different
// renderers for testing without coupling to Lip Gloss.
//
// Width=0 triggers auto-sizing based on content width, which is the recommended
// default to prevent text wrapping in status banners. BorderStyle maps to Lip Gloss
// border presets but is stored as a string to avoid coupling this package to the
// rendering library.
type Box struct {
	Type     BoxType
	Title    string
	Subtitle string
	KVPairs  []KV
	Footer   string

	Width       int
	BorderStyle string
}

// HasContent determines if the box contains any displayable data beyond just
// type information. Used by validators to fail-fast when users attempt to render
// an effectively empty box, which likely indicates a CLI usage error.
func (b *Box) HasContent() bool {
	return b.Title != "" || b.Subtitle != "" || len(b.KVPairs) > 0 || b.Footer != ""
}
