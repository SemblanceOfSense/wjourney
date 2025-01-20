package textgeneration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type aiRequest struct {
    Model string `json:"model"`
    Messages []message `json:"messages"`
}

type message struct {
    Role string `json:"role"`
    Content string `json:"content"`
}

type response struct {
    Id string
    Object string
    Created int
    Model string
    Choices []choice
    Usage usage
    system_fingerprint string
}

type usage struct {
    Prompt_tokens int
    completion_tokens int
    total_tokens int
}

type choice struct {
    Index int
    Message message
    Logprobs string
    Finish_reason string
}

func GetGeneratedText(prompt string, aikey string) (string, error) {
    requestPrompt := message{Role: "user", Content: prompt}
    requestBody := aiRequest{Model: "grok-2-vision-1212", Messages: []message{requestPrompt}}

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", err
    }

    url := "https://api.x.ai/v1/chat/completions"
    req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + aikey)
    if err != nil {
        return "", err
    }

    client := http.Client{Timeout: 10 * time.Second}
    res, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()
    bodyStruct := &response{}
    body, err := io.ReadAll(res.Body)
    err = json.Unmarshal(body, bodyStruct)
    if err != nil { return "", err }

    return bodyStruct.Choices[0].Message.Content, nil
}
