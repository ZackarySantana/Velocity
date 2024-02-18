// Code generated from Pkl module `request_prebuilt`. DO NOT EDIT.
package method

import (
	"encoding"
	"fmt"
)

type Method string

const (
	Get    Method = "get"
	Post   Method = "post"
	Put    Method = "put"
	Delete Method = "delete"
)

// String returns the string representation of Method
func (rcv Method) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(Method)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Method.
func (rcv *Method) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "get":
		*rcv = Get
	case "post":
		*rcv = Post
	case "put":
		*rcv = Put
	case "delete":
		*rcv = Delete
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid Method`, str)
	}
	return nil
}
