package main

import (
	"testing"
)

func TestGetBlockSnfts(t *testing.T) {
	for i := 0; i < 10; i++ {
		GetBlockSnfts( "00000000036780")
		GetBlockSnfts( "36780000000000000000")
	}
}