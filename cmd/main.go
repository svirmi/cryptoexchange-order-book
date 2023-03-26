package main

import (
	"context"
	"order-book/pkg/bybit"
)

func main() {
	source := bybit.NewSource()
	source.Start(context.Background())
}
