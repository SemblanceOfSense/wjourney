package imagegeneration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type openAiRequest struct {
    Model string `json:"model"`
    Prompt string `json:"prompt"`
    NumberOfImages int `json:"n"` // Note: with dall-e-3 only one image can be generated at a time
    Size string `json:"size"`
}

type response struct {
    Created int `json:"created"`
    Data []data `json:"data"`
}

type data struct {
    Url string `json:"url"`
}

func GetImageUrl(prompt string, openaikey string) (string, error) {
    requestBody := openAiRequest{
        Model: "dall-e-3",
        Prompt: prompt,
        NumberOfImages: 1,
        Size: "1024x1024",
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil { fmt.Println("Image"); return "", err }

    url := "https://api.openai.com/v1/images/generations"
    req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
    if err != nil { fmt.Println("Image"); return "", err }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + openaikey)

    client := http.Client{Timeout: 20 * time.Second}
    res, err := client.Do(req)
    if err != nil { fmt.Println("Image"); return "", err }
    defer res.Body.Close()

    bodyStruct := &response{}
    body, err := io.ReadAll(res.Body)
    if err != nil { fmt.Println("Image"); return "", err }
    err = json.Unmarshal(body, bodyStruct)
    if err != nil { fmt.Println("Image"); return "", err }

    if len(bodyStruct.Data) == 0 { return "rejected", nil }
    return bodyStruct.Data[0].Url, nil
}
