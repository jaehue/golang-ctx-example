package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const jobCount = 10

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	quits := make(chan struct{}, 10)

	go func() {
		fmt.Println("Signal:", <-exit)
		for i := 0; i < jobCount; i++ {
			quits <- struct{}{}
		}
	}()

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < jobCount; i++ {
		wg.Add(1)

		go func(quit <-chan struct{}) {
			defer wg.Done()
			result, err := longFuncWithoutCtx(quit)
			fmt.Printf("duration:%v result:%s\n", time.Since(start), result)
			if err != nil {
				fmt.Println(err)
			}
		}(quits)
	}

	wg.Wait()
}

func longFuncWithoutCtx(quit <-chan struct{}) (string, error) {
	done := make(chan string)

	go func() {
		done <- longFunc()
	}()

	select {
	case <-quit:
		return "Fail", errors.New("Force quit")
	case result := <-done:
		return result, nil
	}
}
func longFunc() string {
	<-time.After(time.Second * 3)
	return "Success"
}
