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

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	response, _ := io.ReadAll(res.Body)
	fmt.Println(string(response))

}
