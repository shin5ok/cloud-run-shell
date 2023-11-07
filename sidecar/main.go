package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"log/slog"

	"encoding/json"

	"cloud.google.com/go/compute/metadata"
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
LOOP:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Receiving signal...")
			break LOOP
		default:
			data, err := f()
			if err != nil {
				logger.Error(err.Error())
				return
			}
			logger.Info(
				"instance_id="+instanceId,
				"data", data,
				"severity", "info",
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

func getCgroups() (string, error) {
	result := map[string]string{}
	path := path.Join("/sys/fs/cgroup", "cpuset")

	dir, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	pattern := regexp.MustCompile(`^cpuset\.+`)

	for n, file := range files {

		if !pattern.Match([]byte(file.Name())) {
			continue
		}

		f, err := os.Open(path + "/" + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		data, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dataString := string(data)
		dataString = strings.ReplaceAll(dataString, "\n", ",")

		key := fmt.Sprintf("%d:%s", n, file.Name())
		result[key] = dataString

		f.Close()

	}

	dataJson, _ := json.Marshal(result)

	return string(dataJson), nil
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

	var wg sync.WaitGroup
	wg.Add(1)
	go vlogging(ctx, getCPU, 2)
	go vlogging(ctx, getConnNumber, 5)
	go vlogging(ctx, getCgroups, 5)
	wg.Wait()
}
