// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Locale represents a Unicode locale identifier.
type Locale struct {
	Language  string
	Script    string
	Territory string
}

// NewLocale creates a new Locale from its string representation.
func NewLocale(id string) Locale {
	// Normalize the ID ("SR_rs_LATN" => "sr-Latn-RS").
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, "_", "-")
	locale := Locale{}
	for i, part := range strings.Split(id, "-") {
		if i == 0 {
			locale.Language = part
			continue
		}
		partLen := len(part)
		if partLen == 4 {
			// Uppercase the first letter in a UTF8-safe manner.
			r, size := utf8.DecodeRuneInString(part)
			locale.Script = string(unicode.ToTitle(r)) + part[size:]
			continue
		}
		if partLen == 2 || partLen == 3 {
			locale.Territory = strings.ToUpper(part)
			continue
		}
	}

	return locale
}

// String returns the string representation of l.
func (l Locale) String() string {
	b := strings.Builder{}
	b.WriteString(l.Language)
	if l.Script != "" {
		b.WriteString("-")
		b.WriteString(l.Script)
	}
	if l.Territory != "" {
		b.WriteString("-")
		b.WriteString(l.Territory)
	}

	return b.String()
}

// MarshalText implements the encoding.TextMarshaler interface.
func (l Locale) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (l *Locale) UnmarshalText(b []byte) error {
	*l = NewLocale(string(b))
	return nil
}

// IsEmpty returns whether l is empty.
func (l Locale) IsEmpty() bool {
	return l.Language == "" && l.Script == "" && l.Territory == ""
}
