// Code generated from Pkl module `request_prebuilt`. DO NOT EDIT.
package request

import "github.com/zackarysantana/velocity/gen/pkl/prebuilts/request/method"

type Request struct {
	Method method.Method `pkl:"method"`

	Url string `pkl:"url"`

	Body *string `pkl:"body"`

	Timeout *int `pkl:"timeout"`

	Headers *map[string]string `pkl:"headers"`
}
