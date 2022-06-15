package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestHomePage(t *testing.T) {
	//testCount := 10
	//tKey, tokens, err := Mlogin(testCount)
	//if err != nil {
	//	fmt.Println("HomePage err=", err)
	//	return
	//}
	//fmt.Println("HomePage end.")
	//fmt.Println("HomePage Test NewCollections.")
	spendT := time.Now()
	fmt.Printf("HomePage() QueryHomePage i=%-4d spend time=%s time.now=%s\n", 10, time.Now().Sub(spendT), time.Now())
	testCount := 1000
	wd := sync.WaitGroup{}
	for i := 0; i < testCount; i++ {
		wd.Add(1)
		go func(i int) {
			defer wd.Done()
			spendT := time.Now()
			err := QueryHomePage()
			if err != nil {
				fmt.Println("HomePage() err=", err)
			}
			fmt.Printf("HomePage() QueryHomePage i=%-4d spend time=%-10s time.now=%s\n", i, time.Now().Sub(spendT), time.Now())
		}(i)
	}
	wd.Wait()

	fmt.Println("end test HomePage.")
}