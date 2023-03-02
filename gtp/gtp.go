package gtp

import (
	"bytes"
	"encoding/json"
	"github.com/869413421/wechatbot/config"
	"io/ioutil"
	"log"
	"net/http"
)

const BASEURL = "https://api.openai.com/v1/chat/"


// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      		string                   `json:"id"`
	Object  		string                   `json:"object"`
	Created 		int                      `json:"created"`
	Model   		string                   `json:"model"`
	// Choices 		[]map[string]map[string]interface{} `json:"choices"`

	Choices []struct {
		Message struct {
			Role     string                   `json:"role"`
			Content  string                   `json:"content"`
		} `json:"message"`
	} `json:"choices"`

	Usage   		map[string]interface{}   `json:"usage"`
	FinishReason 	string 		     		 `json:"finish_reason"`
	Index			int64					 `json:index`
}

type ChoiceItem struct {
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  				   `json:"model"`
	Messages         []map[string]interface{}  `json:"messages"`
	MaxTokens        int     				   `json:"max_tokens"`
	Temperature      float32                   `json:"temperature"`
	TopP             int                       `json:"top_p"`
	FrequencyPenalty int                       `json:"frequency_penalty"`
	PresencePenalty  int                       `json:"presence_penalty"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
	message := []map[string]interface{}{
		{"role": "system", "content": "You are a useful assistant."},
		{"role": "user", "content": msg},
	}

	requestBody := ChatGPTRequestBody{
		Model:            "gpt-3.5-turbo",
		Messages:         message,
		MaxTokens:        3000,
		Temperature:      0,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// log.Printf("gpt response text: %s \n", body["choices"][0]["message"]["content"])
	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	
	log.Printf("gpt response text: %s \n", gptResponseBody.Choices[0].Message.Content)
	if err != nil {
		return "", err
	}
	var reply string
	// if len(gptResponseBody.Choices) > 0 {
	// 	for _, v := range gptResponseBody.Choices {
	// 		reply = v["message"]["content"].(string)
	// 		break
	// 	}
	// }
	reply = gptResponseBody.Choices[0].Message.Content

	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
