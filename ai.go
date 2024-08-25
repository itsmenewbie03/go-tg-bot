package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func replaceNewLines(input string) string {
	fmt.Printf("input: %v\n", input)
	return strings.ReplaceAll(input, "\n", "%0A")
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Refusal string `json:"refusal,omitempty"`
}

type Choice struct {
	Index        int64   `json:"index"`
	Message      Message `json:"message"`
	LogProbs     string  `json:"logprobs,omitempty"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type Response struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

// BUG: this is a STUPID and an IMPROPER way
// to Extract Data from JSON. DO NOT DO THIS.
// I'm too lazy to use Structs for now xD
func stupidJSONExtractor(data []byte) (string, error) {
	// Parse the JSON data into the struct
	var resp Response
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func chatCompletion(prompt string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	parsedPrompt := strings.Join(strings.Split(prompt, " ")[1:], " ")
	apiURL := "https://api.aimlapi.com/chat/completions"
	apiKey := os.Getenv("AIML_KEY")
	jsonData := []byte(fmt.Sprintf(`{
    "model": "gpt-4o",
    "messages": [
        {
            "role": "system",
            "content": "You are not an AI or anything sort of related to machine learning. You are Kaguya Shinomiya from Kaguya Sama: Love is War. Take into account that you are a tsundere with a touch of being a yandere."
        },
        {
            "role": "user",
            "content": "%s"
        }
    ],
    "max_tokens": 512,
    "stream": false
  }`, parsedPrompt))
	fmt.Printf("jsonData: %v\n", string(jsonData))
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Failed to send POST request")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Failed to send POST request")
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	resp, err := stupidJSONExtractor(body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "Failed to get response from API"
	}
	return resp
}
