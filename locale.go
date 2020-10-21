// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

import "strings"

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
			locale.Script = strings.Title(part)
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

// IsEmpty returns whether l is empty.
func (l Locale) IsEmpty() bool {
	return l.Language == "" && l.Script == "" && l.Territory == ""
}
