// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bojanz/address"
)

// testFormat is a reduced format for testing purposes.
type testFormat struct {
	Locale     string            `json:"locale"`
	Layout     string            `json:"layout"`
	RegionType string            `json:"region_type"`
	Regions    map[string]string `json:"regions"`
}

func TestFormatHandlerNoLocale(t *testing.T) {
	req, err := http.NewRequest("GET", "/address-formats", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := address.FormatHandler{}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("got HTTP %v want HTTP %v", status, http.StatusOK)
	}
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("got %v want %v", contentType, "application/json")
	}
	if contentLanguage := rr.Header().Get("Content-Language"); contentLanguage != "en" {
		t.Errorf("got %v want %v", contentLanguage, "en")
	}

	var data map[string]testFormat
	err = json.Unmarshal(rr.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}
	format, ok := data["TW"]
	wantFormat := address.GetFormat("TW")
	if !ok {
		t.Errorf(`address format "TW" not found.`)
	}
	// Confirm that Locale and RegionType were correctly converted to strings.
	if format.Locale != wantFormat.Locale.String() {
		t.Errorf("got %v, want %v", format.Locale, wantFormat.Locale.String())
	}
	if format.RegionType != wantFormat.RegionType.String() {
		t.Errorf("got %v, want %v", format.Layout, wantFormat.Layout)
	}
	// Confirm that the correct layout and regions were selected.
	if format.Layout != wantFormat.Layout {
		t.Errorf("got %q, want %q", format.Layout, wantFormat.Layout)
	}
	if !reflect.DeepEqual(format.Regions, wantFormat.Regions) {
		t.Errorf("got %v, want %v", format.Regions, wantFormat.Regions)
	}
}

func TestFormatHandlerLocaleQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/address-formats?locale=zh", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := address.FormatHandler{}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("got HTTP %v want HTTP %v", status, http.StatusOK)
	}
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("got %v want %v", contentType, "application/json")
	}
	if contentLanguage := rr.Header().Get("Content-Language"); contentLanguage != "zh" {
		t.Errorf("got %v want %v", contentLanguage, "zh")
	}

	var data map[string]testFormat
	err = json.Unmarshal(rr.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}
	format, ok := data["TW"]
	wantFormat := address.GetFormat("TW")
	if !ok {
		t.Errorf(`address format "TW" not found.`)
	}
	// Confirm that the correct layout and regions were selected.
	if format.Layout != wantFormat.LocalLayout {
		t.Errorf("got %q, want %q", format.Layout, wantFormat.LocalLayout)
	}
	if !reflect.DeepEqual(format.Regions, wantFormat.LocalRegions) {
		t.Errorf("got %v, want %v", format.Regions, wantFormat.LocalRegions)
	}
}

func TestFormatHandlerLocaleHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/address-formats", nil)
	req.Header.Add("Accept-Language", "zh-hant, zh, en;q=0.8")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := address.FormatHandler{}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("got HTTP %v want HTTP %v", status, http.StatusOK)
	}
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("got %v want %v", contentType, "application/json")
	}
	if contentLanguage := rr.Header().Get("Content-Language"); contentLanguage != "zh-Hant" {
		t.Errorf("got %v want %v", contentLanguage, "zh-Hant")
	}

	var data map[string]testFormat
	err = json.Unmarshal(rr.Body.Bytes(), &data)
	if err != nil {
		t.Fatal(err)
	}
	format, ok := data["TW"]
	wantFormat := address.GetFormat("TW")
	if !ok {
		t.Errorf(`address format "TW" not found.`)
	}
	// Confirm that the correct layout and regions were selected.
	if format.Layout != wantFormat.LocalLayout {
		t.Errorf("got %q, want %q", format.Layout, wantFormat.LocalLayout)
	}
	if !reflect.DeepEqual(format.Regions, wantFormat.LocalRegions) {
		t.Errorf("got %v, want %v", format.Regions, wantFormat.LocalRegions)
	}
}
