package main

import (
	"testing"
)

func TestGetBlockTrans(t *testing.T) {
	GetBlockTrans("68450")
	GetBlockTrans("68528")
	GetBlockTrans("69881")
	GetBlockTrans("81987")
	GetBlockTrans("84696")
}
