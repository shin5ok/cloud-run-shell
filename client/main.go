package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Commanding struct {
	Command string `json:"command"`
}

type Responsing struct {
	Stdout     string `json:"stdout_output"`
	Stderr     string `json:"stderr_output"`
	ReturnCode int    `json:"return_code"`
}

var (
	secret     = os.Getenv("SECRET")
	authHeader = "X-MyGCP-Secret"
	url        = os.Getenv("URL")
	token      = os.Getenv("TOKEN")
)

func main() {
	cmd := Commanding{Command: strings.Join(os.Args[1:], " ")}
	data, _ := json.Marshal(cmd)
	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url+"/shellcommand",
		bytes.NewBuffer(data),
	)

	req.Header.Add(authHeader, secret)
	if err != nil {
		log.Fatal(err, req)
	}

	if token != "" {
		req.Header.Add(`Authorization`, "Bearer "+token)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(res.Status, err)
	}
	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(response))
}
