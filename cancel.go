package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Println("Signal:", <-exit)
		cancel()
	}()

	start := time.Now()
	result, err := longFuncWithCtx(ctx)
	fmt.Printf("duration:%v result:%s\n", time.Since(start), result)
	if err != nil {
		log.Fatal(err)
	}

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
