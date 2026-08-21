// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

import (
	"html"
	"strings"
)

// Formatter formats addresses for display.
type Formatter struct {
	locale Locale
	// CountryMapper maps country codes to country names.
	// Can be used to retrieve country names from another (localized) source.
	// Defaults to a function that uses English country names included in the package.
	CountryMapper func(countryCode string, locale Locale) string
	// NoCountry turns off displaying the country name.
	// Defaults to false.
	NoCountry bool
	// WrapperElement is the wrapper HTML element.
	// Defaults to "p".
	WrapperElement string
	// WrapperClass is the wrapper HTML class.
	// Defaults to "address".
	WrapperClass string
}

// NewFormatter creates a new formatter for the given locale.
func NewFormatter(locale Locale) *Formatter {
	f := &Formatter{
		locale: locale,
		CountryMapper: func(countryCode string, locale Locale) string {
			return countries[countryCode]
		},
		WrapperElement: "p",
		WrapperClass:   "address",
	}
	return f
}

// Locale returns the locale.
func (f *Formatter) Locale() Locale {
	return f.locale
}

// Format formats the given address.
func (f *Formatter) Format(addr Address) string {
	if addr.IsEmpty() {
		return ""
	}
	format := GetFormat(addr.CountryCode)
	layout := format.SelectLayout(f.locale)
	countryBefore := (layout == format.LocalLayout)
	countryAfter := (layout != format.LocalLayout)
	country := ""
	if !f.NoCountry {
		country = html.EscapeString(f.CountryMapper(addr.CountryCode, f.locale))
		country = `<span class="country" data-value="` + addr.CountryCode + `">` + country + `</span>`
	}
	values := f.getValues(addr)
	for field, value := range values {
		if value != "" {
			value = html.EscapeString(value)
			value = `<span class="` + f.getClass(field) + `">` + value + `</span>`
			values[field] = value
		}
	}

	sb := strings.Builder{}
	sb.Grow(200)
	sb.WriteString(`<` + f.WrapperElement + ` class="`)
	sb.WriteString(f.WrapperClass)
	sb.WriteString(`" translate="no">` + "\n")
	if !f.NoCountry && countryBefore {
		sb.WriteString(country)
		sb.WriteString("<br>\n")
	}
	f.writeValues(&sb, layout, values)
	if !f.NoCountry && countryAfter {
		sb.WriteString("<br>\n")
		sb.WriteString(country)
	}
	sb.WriteString("\n</" + f.WrapperElement + ">")

	return sb.String()
}

// getClass returns the HTML class for the given field.
func (f *Formatter) getClass(field Field) string {
	var class string
	switch field {
	case FieldLine1:
		class = "line1"
	case FieldLine2:
		class = "line2"
	case FieldLine3:
		class = "line3"
	case FieldSublocality:
		class = "sublocality"
	case FieldLocality:
		class = "locality"
	case FieldRegion:
		class = "region"
	case FieldPostalCode:
		class = "postal-code"
	}

	return class
}

// getValues returns all values for the given address, keyed by field.
//
// Region IDs are replaced by region names if available.
func (f *Formatter) getValues(addr Address) map[Field]string {
	values := map[Field]string{
		FieldLine1:       addr.Line1,
		FieldLine2:       addr.Line2,
		FieldLine3:       addr.Line3,
		FieldSublocality: addr.Sublocality,
		FieldLocality:    addr.Locality,
		FieldRegion:      addr.Region,
		FieldPostalCode:  addr.PostalCode,
	}
	format := GetFormat(addr.CountryCode)
	regions := format.SelectRegions(f.locale)
	if !format.ShowRegionID && regions.Len() > 0 {
		region, ok := regions.Get(addr.Region)
		if ok {
			values[FieldRegion] = region
		}
	}

	return values
}

// writeValues writes the formatted address one line at a time, skipping any
// lines that have no values.
func (f *Formatter) writeValues(b *strings.Builder, layout string, values map[Field]string) {
	written := false
	for len(layout) > 0 {
		line := layout
		if n := strings.IndexByte(layout, '\n'); n >= 0 {
			line, layout = layout[:n], layout[n+1:]
		} else {
			layout = ""
		}
		if !hasValues(line, values) {
			continue
		}
		if written {
			b.WriteString("<br>\n")
		}
		writeLine(b, line, values)
		written = true
	}
}

// hasValues reports whether a line of the layout has at least one value.
func hasValues(line string, values map[Field]string) bool {
	for i := 0; i+1 < len(line); i++ {
		if line[i] == '%' && values[Field(line[i+1:i+2])] != "" {
			return true
		}
	}

	return false
}

// writeLine fills in and writes a single line of the address layout.
//
// Empty fields are left out, along with any separator between them and the
// next field that has a value.
func writeLine(b *strings.Builder, line string, values map[Field]string) {
	prev := 0
	sep := ""
	written := false
	skipped := false
	for i := 0; i+1 < len(line); i++ {
		if line[i] != '%' {
			continue
		}
		field := Field(line[i+1 : i+2])
		if !skipped {
			// Save the separator before the first skipped field. Any separators
			// after it belong to fields that are being left out.
			sep = line[prev:i]
		}
		prev = i + 2
		if values[field] == "" {
			skipped = true
			continue
		}
		if written || !skipped {
			b.WriteString(sep)
		}
		b.WriteString(values[field])
		written = true
		skipped = false
	}
}
