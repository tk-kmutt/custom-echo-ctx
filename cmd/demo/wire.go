//+build wireinject

package main

import (
	"custom-echo-ctx/internal"

	"github.com/google/wire"
)

func Ready() *internal.Server {
	wire.Build(
		internal.NewServer,
	)

	return &internal.Server{}
}

func MockReady() *internal.Server {
	wire.Build(
		internal.NewServer,
	)

	return &internal.Server{}
}
