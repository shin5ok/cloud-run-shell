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
	"time"

	"github.com/briandowns/spinner"
)

type Commanding struct {
	Command string `json:"command"`
}

type ResponseData struct {
	Stdout     []string `json:"stdout_output"`
	Stderr     []string `json:"stderr_output"`
	Message    string   `json:"message"`
	ReturnCode int      `json:"return_code"`
}

var (
	authHeader = "X-MyGCP-Secret"
	secret     = os.Getenv("SECRET")
	url        = os.Getenv("URL")
	token      = os.Getenv("TOKEN")
	modeJson   = os.Getenv("MODE_JSON")
)

func init() {
	if url == "" {
		log.Fatal(`
	Set url like this,
	export URL=https://xxxxxxx-an.a.run.app
		`)
	}
}

func main() {

	s := spinner.New(spinner.CharSets[9], 300*time.Millisecond)
	s.Start()

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

	s.Stop()

	if modeJson != "" {
		fmt.Println(string(response))
		return
	}

	var r ResponseData
	if err := json.Unmarshal(response, &r); err != nil {
		log.Panic(err)
	}

	if r.ReturnCode != 0 {
		for _, line := range r.Stderr {
			fmt.Println(line)
		}
	}

	for _, line := range r.Stdout {
		fmt.Println(line)
	}

}
