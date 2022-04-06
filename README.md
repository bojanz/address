# address [![Build](https://github.com/bojanz/address/actions/workflows/build.yml/badge.svg)](https://github.com/bojanz/address/actions/workflows/build.yml) [![Coverage Status](https://coveralls.io/repos/github/bojanz/address/badge.svg?branch=master)](https://coveralls.io/github/bojanz/address?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/bojanz/address)](https://goreportcard.com/report/github.com/bojanz/address) [![PkgGoDev](https://pkg.go.dev/badge/github.com/bojanz/address)](https://pkg.go.dev/github.com/bojanz/address)

Handles address representation, validation and formatting.

Inspired by Google's [libaddressinput](https://github.com/google/libaddressinput).

Backstory: https://bojanz.github.io/address-handling-go/

## Features

1. Address struct.
2. Address formats for ~200 countries.
3. Regions for ~50 countries, with local names where relevant (e.g: Okinawa / 沖縄県).
4. Country list, powered by CLDR v41.
5. HTML formatter.
6. HTTP handler for serving address formats and regions as JSON: only ~14kb gzipped!

## Address struct

Represents an address as commonly handled by web applications and APIs.

```go
type Address struct {
	Line1 string
	Line2 string
	Line3 string
	// Sublocality is the neighborhood/suburb/district.
	Sublocality string
	// Locality is the city/village/post town.
	Locality string
	// Region is the state/province/prefecture.
	// An ISO code is used when available.
	Region string
	// PostalCode is the postal/zip/pin code.
	PostalCode string
	// CountryCode is the two-letter code as defined by CLDR.
	CountryCode string
}
```

Recipient fields such as FirstName/LastName are not included since they are usually
present on the parent (Contact/Customer/User) struct. This allows the package
to avoid being opinionated about name handling.

There are three line fields, matching the HTML5 autocomplete spec and many shipping APIs.
This leaves enough space for specifying an organization and department, house or hotel
name, and other similar "care of" use cases. When mapping to an API that only has two address lines,
Line3 can be appended to Line2, separated by a comma.

## Address formats

The following information [is available](https://github.com/bojanz/address/blob/master/formats.go#L6):

- Which fields are used, and in which order.
- Which fields are required.
- Labels for the sublocality, locality, region and postal code fields.
- Regular expression pattern for validating postal codes.
- Regions and how to display them in an address.

Certain countries (e.g. China, Japan, Russia, Ukraine) have region names defined in both Latin and local scripts. The script is selected based on locale. For example, the "ru" locale will use Russian regions in Cyrilic, while "ru-Latn" and other locales will use the Latin version.

[Helpers](https://github.com/bojanz/address/blob/master/address.go#L61) are provided for validating required fields, regions, postal codes.

Format data was generated from Google's [Address Data](https://chromium-i18n.appspot.com/ssl-address) but isn't
automatically regenerated, to allow the community to submit their own corrections directly to the package.

## Countries

The country list is auto-generated from CLDR. 
Updating to the latest CLDR release is always one `go generate` away.

Most software uses the CLDR country list instead of the ISO one because the CLDR country names match their colloquial usage more closely (e.g. "Russia" instead of "Russian Federation"). 

To reduce the size of the included data, this package only includes country names in English.
Translated country names can be fetched on the frontend via [Intl.DisplayNames](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DisplayNames). Alternatively, one can plug in [x/text/language/display](https://pkg.go.dev/golang.org/x/text/language/display) by setting a custom CountryMapper on the formatter.

## Formatter

Displays an address as HTML, using the country's address format.

The wrapper element ("p") and class ("address") can be configured.
The country name can be omitted, for the use case where all addresses belong to the same country. 

```go
addr := address.Address{
    Line1:       "1098 Alta Ave",
    Locality:    "Mountain View",
    Region:      "CA",
    PostalCode:  "94043",
    CountryCode: "US",
}
locale := address.NewLocale("en")
formatter := address.NewFormatter(locale)
output := formatter.Format(addr)
// Output:
// <p class="address" translate="no">
// <span class="line1">1098 Alta Ave</span><br>
// <span class="locality">Mountain View</span>, <span class="region">CA</span> <span class="postal-code">94043</span><br>
// <span class="country" data-value="US">United States</span>
// </p>

addr = address.Address{
    Line1:       "幸福中路",
    Sublocality: "新城区",
    Locality:    "西安市",
    Region:      "SN",
    PostalCode:  "710043",
    CountryCode: "CN",
}
locale := address.NewLocale("zh")
formatter := address.NewFormatter(locale)
formatter.NoCountry = true
formatter.WrapperElement = "div"
formatter.WrapperClass = "postal-address"
output := formatter.Format(addr)
// Output:
// <div class="postal-address" translate="no">
// <span class="postal-code">710043</span><br>
// <span class="region">陕西省</span><span class="locality">西安市</span><span class="sublocality">新城区</span><br>
// <span class="line1">幸福中路</span>
// </div>
```
