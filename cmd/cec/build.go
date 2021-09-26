// +build !mock

package main

import "custom-echo-ctx/internal"

func NewServer() *internal.Server {
	return Ready()
}
