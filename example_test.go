// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address_test

import (
	"fmt"

	"github.com/bojanz/address"
)

func ExampleNewLocale() {
	firstLocale := address.NewLocale("en-US")
	fmt.Println(firstLocale)
	fmt.Println(firstLocale.Language, firstLocale.Territory)

	// Locale IDs are normalized.
	secondLocale := address.NewLocale("sr_rs_latn")
	fmt.Println(secondLocale)
	fmt.Println(secondLocale.Language, secondLocale.Script, secondLocale.Territory)
	// Output: en-US
	// en US
	// sr-Latn-RS
	// sr Latn RS
}

func ExampleFormatter_Format() {
	locale := address.NewLocale("en")
	formatter := address.NewFormatter(locale)
	addr := address.Address{
		Line1:       "1098 Alta Ave",
		Locality:    "Mountain View",
		Region:      "CA",
		PostalCode:  "94043",
		CountryCode: "US",
	}
	fmt.Println(formatter.Format(addr))

	addr = address.Address{
		Line1:       "幸福中路",
		Sublocality: "新城区",
		Locality:    "西安市",
		Region:      "SN",
		PostalCode:  "710043",
		CountryCode: "CN",
	}
	locale = address.NewLocale("zh")
	formatter = address.NewFormatter(locale)
	formatter.NoCountry = true
	formatter.WrapperElement = "div"
	formatter.WrapperClass = "postal-address"
	fmt.Println(formatter.Format(addr))
	// Output:
	// <p class="address" translate="no">
	// <span class="line1">1098 Alta Ave</span><br>
	// <span class="locality">Mountain View</span>, <span class="region">CA</span> <span class="postal-code">94043</span><br>
	// <span class="country" data-value="US">United States</span>
	// </p>
	// <div class="postal-address" translate="no">
	// <span class="postal-code">710043</span><br>
	// <span class="region">陕西省</span><span class="locality">西安市</span><span class="sublocality">新城区</span><br>
	// <span class="line1">幸福中路</span>
	// </div>
}
