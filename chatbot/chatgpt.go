package chatbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func fetchSupabaseData() ([]byte, error) {
	url := fmt.Sprintf("%s/table_name", os.Getenv("SUPABASE_URL")) // Replace 'table_name' with your table name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", os.Getenv("SUPABASE_API_KEY"))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func generateChatGPTResponse(prompt string, contextData []byte) (string, error) {
	type ChatGPTRequest struct {
		Prompt   string `json:"prompt"`
		MaxTurns int    `json:"max_turns"`
	}

	type ChatGPTResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	chatGPTReq := ChatGPTRequest{
		Prompt:   prompt,
		MaxTurns: 5, // Number of conversation turns
	}
	if len(contextData) > 0 {
		chatGPTReq.Prompt += "\n\n" + string(contextData)
	}

	reqBody, err := json.Marshal(chatGPTReq)
	if err != nil {
		return "", err
	}

	url := os.Getenv("CHAT_GPT_COMPLETION_URL")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CHAT_GPT_API_KEY")))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var chatGPTResp ChatGPTResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(respBody, &chatGPTResp); err != nil {
		return "", err
	}

	return chatGPTResp.Choices[0].Message.Content, nil
}

func AutoCompletions(userInput string) (string, error) {
	// Fetch data from Supabase
	supabaseData, err := fetchSupabaseData()
	if err != nil {
		fmt.Println("Error fetching data from Supabase:", err)
		return "", err
	}

	// Generate response using ChatGPT
	response, err := generateChatGPTResponse(userInput, supabaseData)
	if err != nil {
		fmt.Println("Error generating ChatGPT response:", err)
		return "", err
	}

	return fmt.Sprintf("Chatbot: %s", response), nil
}
