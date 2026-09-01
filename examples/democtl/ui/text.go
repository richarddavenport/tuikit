// Package-level text helpers.
//
// What used to be here — trim, clip, pad, padVisible, rule, wrap, wrapFirst and
// a 20-line composite — is gone. The canvas clips, so nothing needs padding to
// a width or truncating to fit a border, and a box drawn last is on top without
// anything being composited underneath it. What survives is not text handling
// at all: it is two ways of saying how long something took.
package ui

import (
	"fmt"
	"time"
)

// ago is a duration as a person would say it. Coarse on purpose: "4m ago" is
// what the reader wants, and "4m13.204s ago" is the same fact made unreadable.
func ago(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	default:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
}

// took is how long a step ran, at the precision a step list can use.
func took(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
