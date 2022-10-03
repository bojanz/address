// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

import (
	"encoding/json"
	"net/http"
	"strings"
)

// FormatHandler is an HTTP handler for serving address formats.
//
// Response size is ~47kb, or ~14kb if gzip compression is used.
//
// The locale can be provided either as a query string (?locale=fr)
// or as a header (Accept-Language:fr). Defaults to "en".
type FormatHandler struct{}

// ServeHTTP implements the http.Handler interface.
func (h *FormatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	locale := h.getLocale(r)
	// Preselecting the layout and regions reduces HTTP request size by ~20%.
	type localizedFormat struct {
		Locale            string          `json:"locale,omitempty"`
		Layout            string          `json:"layout,omitempty"`
		Required          []Field         `json:"required,omitempty"`
		SublocalityType   SublocalityType `json:"sublocality_type,omitempty"`
		LocalityType      LocalityType    `json:"locality_type,omitempty"`
		RegionType        RegionType      `json:"region_type,omitempty"`
		PostalCodeType    PostalCodeType  `json:"postal_code_type,omitempty"`
		PostalCodePattern string          `json:"postal_code_pattern,omitempty"`
		ShowRegionID      bool            `json:"show_region_id,omitempty"`
		Regions           *RegionMap      `json:"regions,omitempty"`
	}
	data := make(map[string]localizedFormat, len(formats))
	for countryCode, format := range formats {
		lf := localizedFormat{
			Locale:            format.Locale.String(),
			Layout:            format.SelectLayout(locale),
			Required:          format.Required,
			SublocalityType:   format.SublocalityType,
			LocalityType:      format.LocalityType,
			RegionType:        format.RegionType,
			PostalCodeType:    format.PostalCodeType,
			PostalCodePattern: format.PostalCodePattern,
			ShowRegionID:      format.ShowRegionID,
		}
		if regions := format.SelectRegions(locale); regions.Len() > 0 {
			lf.Regions = &regions
		}
		data[countryCode] = lf
	}
	jsonData, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Language", locale.String())
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// getLocale returns the locale to use.
//
// Priority:
// 1) Query string (?locale=fr)
// 2) Header (Accept-Language=fr)
// 3) English
func (h *FormatHandler) getLocale(r *http.Request) Locale {
	var locale Locale
	if param := r.URL.Query().Get("locale"); param != "" {
		locale = NewLocale(param)
	} else if accept := r.Header.Get("Accept-Language"); accept != "" {
		for _, sep := range []string{",", ";"} {
			if strings.Contains(accept, sep) {
				acceptParts := strings.Split(accept, sep)
				accept = acceptParts[0]
			}
		}
		locale = NewLocale(strings.TrimSpace(accept))
	} else {
		locale = Locale{Language: "en"}
	}

	return locale
}
