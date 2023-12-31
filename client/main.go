package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

type Commanding struct {
	Command string `json:"command"`
}

type ResponseData struct {
	Stdout     []string         `json:"stdout_output"`
	Stderr     []string         `json:"stderr_output"`
	Message    string           `json:"message"`
	ReturnCode int              `json:"return_code"`
	Metadata   []map[string]any `json:"metadata"`
}

var (
	authHeader = "X-MyGCP-Secret"
	secret     = os.Getenv("SECRET")
	url        = os.Getenv("URL")
	token      = os.Getenv("TOKEN")
	jsonMode   = os.Getenv("JSON_MODE")
)

func init() {
	if url == "" {
		fmt.Println(`
	Set url like this,
	export URL=https://xxxxxxx-an.a.run.app
		`)
		os.Exit(0)
	}

}

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	s := spinner.New(spinner.CharSets[26], 300*time.Millisecond)
	defer s.Stop()
	s.Start()

	cmd := Commanding{Command: strings.Join(os.Args[1:], " ")}
	data, err := json.Marshal(cmd)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url+"/shellcommand",
		bytes.NewBuffer(data),
	)

	req.Header.Add(authHeader, secret)
	if err != nil {
		panic(err)
	}

	if token != "" {
		req.Header.Add(`Authorization`, fmt.Sprintf("Bearer %s", token))
	}

	go func() {
		<-ctx.Done()
		fmt.Println("exiting with context:", ctx.Err())
		defer stop()
		defer s.Stop()
		os.Exit(1)
	}()

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   1 * time.Hour,
				KeepAlive: 1 * time.Hour,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2: true,
		},
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Something wrong with code: %d...it might have invalid SECRET?", res.StatusCode))
	}

	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	s.Stop()

	if jsonMode != "" {
		fmt.Println(string(response))
		return
	}

	var r ResponseData
	if err := json.Unmarshal(response, &r); err != nil {
		panic(err)
	}

	if r.ReturnCode != 0 {
		for _, line := range r.Stderr {
			fmt.Println(line)
		}
	}

	for _, line := range r.Stdout {
		fmt.Println(line)
	}

	os.Exit(r.ReturnCode)

}
