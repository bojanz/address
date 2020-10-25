// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address_test

import (
	"strings"
	"testing"

	"github.com/bojanz/address"
)

func TestFormatter_Locale(t *testing.T) {
	locale := address.NewLocale("fr-CH")
	formatter := address.NewFormatter(locale)
	got := formatter.Locale().String()
	if got != "fr-CH" {
		t.Errorf("got %v, want fr-CH", got)
	}
}

func TestFormatter_FormatES(t *testing.T) {
	locale := address.NewLocale("en")
	formatter := address.NewFormatter(locale)
	formatter.WrapperElement = "div"
	formatter.WrapperClass = "address postal-address"

	// Empty address.
	got := formatter.Format(address.Address{})
	if got != "" {
		t.Errorf("got: %v, want an empty string", got)
	}

	// Address with embedded HTML and an unrecognized region.
	addr := address.Address{
		Line1:       "Calle Numa<b>55</b>",
		Locality:    "Dos Hermanas",
		Region:      "SEV",
		PostalCode:  "41089",
		CountryCode: "ES",
	}
	wantLines := []string{
		`<div class="address postal-address" translate="no">`,
		`<span class="line1">Calle Numa&lt;b&gt;55&lt;/b&gt;</span><br>`,
		`<span class="postal-code">41089</span> <span class="locality">Dos Hermanas</span> <span class="region">SEV</span><br>`,
		`<span class="country" data-value="ES">Spain</span>`,
		`</div>`,
	}
	got = formatter.Format(addr)
	want := strings.Join(wantLines, "\n")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}

	// Address with no country displayed.
	addr = address.Address{
		Line1:       "Calle Numa 55",
		Locality:    "Dos Hermanas",
		Region:      "SE",
		PostalCode:  "41089",
		CountryCode: "ES",
	}
	wantLines = []string{
		`<div class="address postal-address" translate="no">`,
		`<span class="line1">Calle Numa 55</span><br>`,
		`<span class="postal-code">41089</span> <span class="locality">Dos Hermanas</span> <span class="region">Sevilla</span>`,
		`</div>`,
	}
	formatter.NoCountry = true
	got = formatter.Format(addr)
	want = strings.Join(wantLines, "\n")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}
}

func TestFormatter_FormatUS(t *testing.T) {
	locale := address.NewLocale("en")
	formatter := address.NewFormatter(locale)

	// Full address (every field provided).
	addr := address.Address{
		Line1:       "c/o The Westin Seattle",
		Line2:       "Room #505",
		Line3:       "1900 5th Avenue",
		Locality:    "Seattle",
		Region:      "WA",
		PostalCode:  "98101",
		CountryCode: "US",
	}
	wantLines := []string{
		`<p class="address" translate="no">`,
		`<span class="line1">c/o The Westin Seattle</span><br>`,
		`<span class="line2">Room #505</span><br>`,
		`<span class="line3">1900 5th Avenue</span><br>`,
		`<span class="locality">Seattle</span>, <span class="region">WA</span> <span class="postal-code">98101</span><br>`,
		`<span class="country" data-value="US">United States</span>`,
		`</p>`,
	}
	got := formatter.Format(addr)
	want := strings.Join(wantLines, "\n")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}

	// Partial address (no postal code).
	addr = address.Address{
		Line1:       "1098 Alta Ave",
		Locality:    "Mountain View",
		Region:      "CA",
		CountryCode: "US",
	}
	wantLines = []string{
		`<p class="address" translate="no">`,
		`<span class="line1">1098 Alta Ave</span><br>`,
		`<span class="locality">Mountain View</span>, <span class="region">CA</span><br>`,
		`<span class="country" data-value="US">United States</span>`,
		`</p>`,
	}
	got = formatter.Format(addr)
	want = strings.Join(wantLines, "\n")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}

	// Partial address (no region).
	addr = address.Address{
		Line1:       "1098 Alta Ave",
		Locality:    "Mountain View",
		PostalCode:  "94043",
		CountryCode: "US",
	}
	wantLines = []string{
		`<p class="address" translate="no">`,
		`<span class="line1">1098 Alta Ave</span><br>`,
		`<span class="locality">Mountain View</span> <span class="postal-code">94043</span><br>`,
		`<span class="country" data-value="US">United States</span>`,
		`</p>`,
	}
	got = formatter.Format(addr)
	want = strings.Join(wantLines, "\n")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}
}

func TestFormatter_FormatCN(t *testing.T) {
	locale := address.NewLocale("en")
	formatter := address.NewFormatter(locale)
	// Latin address.
	addr := address.Address{
		Line1:       "Xing Fu Zhong Lu",
		Sublocality: "Xincheng Qu",
		Locality:    "Xi'an Shi",
		Region:      "SN",
		PostalCode:  "710043",
		CountryCode: "CN",
	}
	wantLines := []string{
		`<p class="address" translate="no">`,
		`<span class="line1">Xing Fu Zhong Lu</span><br>`,
		`<span class="sublocality">Xincheng Qu</span><br>`,
		`<span class="locality">Xi&#39;an Shi</span><br>`,
		`<span class="region">Shaanxi Sheng</span>, <span class="postal-code">710043</span><br>`,
		`<span class="country" data-value="CN">China</span>`,
		`</p>`,
	}
	got := formatter.Format(addr)
	want := strings.Join(wantLines, "\n")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}

	locale = address.NewLocale("zh")
	formatter = address.NewFormatter(locale)
	// Local address.
	addr = address.Address{
		Line1:       "幸福中路",
		Sublocality: "新城区",
		Locality:    "西安市",
		Region:      "SN",
		PostalCode:  "710043",
		CountryCode: "CN",
	}
	wantLines = []string{
		`<p class="address" translate="no">`,
		`<span class="country" data-value="CN">China</span><br>`,
		`<span class="postal-code">710043</span><br>`,
		`<span class="region">陕西省</span><span class="locality">西安市</span><span class="sublocality">新城区</span><br>`,
		`<span class="line1">幸福中路</span>`,
		`</p>`,
	}
	got = formatter.Format(addr)
	want = strings.Join(wantLines, "\n")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}
}
