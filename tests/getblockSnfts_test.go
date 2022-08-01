package main

import (
	"strconv"
	"testing"
)

func TestGetBlockSnfts(t *testing.T) {
	for i := 0; i < 10; i++ {
		GetBlockSnfts(strconv.Itoa(i))
	}
}
