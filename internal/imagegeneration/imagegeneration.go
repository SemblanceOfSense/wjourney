package imagegeneration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAiRequest struct {
    Model string `json:"model"`
    Prompt string `json:"prompt"`
    NumberOfImages int `json:"n"` // Note: with dall-e-3 only one image can be generated at a time
    Size string `json:"size"`
}

type Response struct {
    Created int `json:"created"`
    Data []Data `json:"data"`
}

type Data struct {
    Url string `json:"url"`
}

func GetImageUrl(prompt string, openaikey string) (string, error) {
    requestBody := OpenAiRequest{
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

    bodyStruct := &Response{}
    body, err := io.ReadAll(res.Body)
    if err != nil { fmt.Println("Image"); return "", err }
    err = json.Unmarshal(body, bodyStruct)
    if err != nil { fmt.Println("Image"); return "", err }

    return bodyStruct.Data[0].Url, nil
}
