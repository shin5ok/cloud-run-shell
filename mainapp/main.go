package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"log/slog"

	"cloud.google.com/go/compute/metadata"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
var instanceId string
var isDebug bool

func init() {
	orgInstanceName, _ := metadata.InstanceID()
	instanceId = os.Getenv("K_REVISION") + "_" + orgInstanceName[len(orgInstanceName)-6:]
	isDebug = os.Getenv("DEBUG") != ""
}

func vlogging(ctx context.Context, f func() (string, error), waitSecond int) {
	name := os.Getenv("NAME")
LOOP:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Receiving signal...")
			break LOOP
		default:
			data, _ := f()
			logger.Info(
				"instance_id="+instanceId,
				"data", data,
				"severity", "info",
				"name", name,
			)
			time.Sleep(time.Second * time.Duration(waitSecond))
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

func getConnNumber() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "netstat", "--inet", "-n").CombinedOutput()
	if err != nil {
		return err.Error(), err
	}
	re := regexp.MustCompile("ESTABLISHED")
	count := 0
	for n, line := range strings.Split(string(output), "\n") {
		if re.Match([]byte(line)) {
			s := strconv.Itoa(n)
			debugPrint(s, line)
			count += 1
		}
	}
	return fmt.Sprintf("tcp_established: %d", count), nil
}

func debugPrint(line ...string) {
	if isDebug {
		fmt.Println(strings.Join(line, " "))
	}
}

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		stop()
		slog.Info("stopping...", "instance_id", instanceId)
		os.Exit(0)
	}()

	go vlogging(ctx, getCPU, 2)
	go vlogging(ctx, getConnNumber, 5)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":8080", r)
}
