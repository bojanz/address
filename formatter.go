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
	if !format.ShowRegionID && len(regions) > 0 {
		region, ok := regions[addr.Region]
		if ok {
			values[FieldRegion] = region
		}
	}

	return values
}

// writeValues inserts values into the layout and writes it to b.
//
// Tokens of empty fields are removed, as are their preceeding chars.
// For example: "%L, %P" becomes "%L" when %P has no value.
func (f *Formatter) writeValues(b *strings.Builder, layout string, values map[Field]string) {
	prev := 0
	for i := 0; i < len(layout); i++ {
		if layout[i] != '%' {
			continue
		}
		j, k := i+1, i+2
		field := Field(layout[j:k])
		if values[field] != "" {
			prefix := layout[prev:i]
			for l := 0; l < len(prefix); l++ {
				if prefix[l] == '\n' {
					// Prepend <br> to each newline.
					b.WriteString("<br>")
				}
				b.WriteByte(prefix[l])
			}
			b.WriteString(values[field])
		}
		prev = k
	}
}
