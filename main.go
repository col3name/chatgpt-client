package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type ChatGptResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type ChatGptRequest struct {
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
	Prompt    string `json:"prompt"`
}

func main() {
	text := flag.String("t", "write story about sherlock holms", "-t=\"write sherlok holmes story \"")
	token := flag.String("token", "token", "-token={secretToken}")
	flag.Parse()
	if len(*text) == 0 {
		fmt.Println(*text, "-t must be not empty")
		return
	}

	url := "https://api.openai.com/v1/completions"
	method := "POST"

	chatGptRequest := ChatGptRequest{Model: "text-davinci-003", MaxTokens: 2000, Prompt: *text}

	client := &http.Client{}
	marshal, err := json.Marshal(chatGptRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(marshal))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+*token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var respGpt ChatGptResponse
	err = json.Unmarshal(body, &respGpt)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(body))
		return
	}
	for _, choice := range respGpt.Choices {
		fmt.Println(choice.Text)
	}
}
