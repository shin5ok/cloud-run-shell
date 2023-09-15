package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"log/slog"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func vlogging(ctx context.Context, f func() (string, error)) {
	for {
		data, _ := f()
		logger.Info(
			"data", data,
			"severity", "info",
		)
		time.Sleep(time.Duration(time.Second * 1))
	}
}

func getCPU() (string, error) {
	r, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	s := bufio.NewScanner(r)
	re, err := regexp.Compile(`^processor`)
	if err != nil {
		return "", err
	}
	count := 0
	for s.Scan() {
		if re.Match([]byte(s.Text())) {
			count += 1
		}
	}

	data := fmt.Sprintf("cpu_count:%d", count)
	return data, nil
}

func main() {
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go vlogging(ctx, getCPU)
	wg.Wait()
}
