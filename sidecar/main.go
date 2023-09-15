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

	"cloud.google.com/go/compute/metadata"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
var instanceId string

func init() {
	orgInstanceId, _ := metadata.InstanceID()
	instanceId = orgInstanceId[0:10]
}

func vlogging(ctx context.Context, f func() (string, error)) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Receiving signal...")
			break LOOP
		default:
			data, _ := f()
			logger.Info(
				"instance_id: "+instanceId,
				"data", data,
				"severity", "info",
			)
			time.Sleep(time.Duration(time.Second * 1))
		}
	}
	os.Exit(1)
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
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go vlogging(ctx, getCPU)
	wg.Wait()
}
