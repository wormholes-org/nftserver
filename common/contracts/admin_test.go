package contracts

import (
	"fmt"
	//"github.com/nftexchange/nftserver/models"
	"testing"
)

func TestGetAdminList(t *testing.T) {
	EthNode = "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	AdminAddr = "0x56c971ebBC0cD7Ba1f977340140297C0B48b7955"
	addr, err := AdminList()
	if err != nil {
		fmt.Println("get admin error.")
	}
	fmt.Println(addr)
}