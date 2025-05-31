// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package txtutil_test

import (
	"fmt"

	"github.com/aileron-gateway/aileron-gateway/kernel/txtutil"
)

func ExampleTemplate_Content_one() {
	tplString := `this is the text template`
	tpl, err := txtutil.NewTemplate(txtutil.TplText, tplString)
	if err != nil {
		panic("handle error here")
	}

	b := tpl.Content(nil)
	fmt.Println(string(b))
	// Output:
	// this is the text template
}

func ExampleTemplate_Content_two() {
	tplString := `this is the {{.name}} template`
	tpl, err := txtutil.NewTemplate(txtutil.TplGoText, tplString)
	if err != nil {
		panic("handle error here")
	}

	info := map[string]any{
		"name": "Go Text",
		"num":  20,
	}
	b := tpl.Content(info)
	fmt.Println(string(b))
	// Output:
	// this is the Go Text template
}

func ExampleTemplate_Content_three() {
	tplString := `this is the {{.name}} template`
	tpl, err := txtutil.NewTemplate(txtutil.TplGoHTML, tplString)
	if err != nil {
		panic("handle error here")
	}

	info := map[string]any{
		"name": "Go HTML",
		"num":  20,
	}
	b := tpl.Content(info)
	fmt.Println(string(b))
	// Output:
	// this is the Go HTML template
}
