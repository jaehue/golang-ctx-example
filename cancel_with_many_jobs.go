package main

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Println("Signal:", <-exit)
		cancel()
	}()

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < jobCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			result, err := longFuncWithCtx(ctx)
			fmt.Printf("duration:%v result:%s\n", time.Since(start), result)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	wg.Wait()

}

func longFuncWithCtx(ctx context.Context) (string, error) {
	done := make(chan string)

	go func() {
		done <- longFunc()
	}()

	select {
	case <-ctx.Done():
		return "Fail", ctx.Err()
	case result := <-done:
		return result, nil
	}
}

func longFunc() string {
	<-time.After(time.Second * 3)
	return "Success"
}
