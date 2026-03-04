package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenRouterClient struct {
	apiKey string
	model  string
	url    string
}

func NewOpenRouterClient(apiKey, model, url string) *OpenRouterClient {
	return &OpenRouterClient{apiKey: apiKey, 
		model: model, 
		url: url,
	}
}

func (c *OpenRouterClient) GenerateSemanticTags(name, description string) (string, error) {
	prompt := fmt.Sprintf(
		"Analyze this food item: Name: %s, Description: %s. "+
			"Generate a short list of semantic search keywords including flavor profiles, dietary info, and mood. "+
			"Output ONLY the keywords separated by commas.", name, description)

	payload := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", c.url+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "http://localhost:8084")
	req.Header.Set("X-Title", "Go-Food-Delivery-Project")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		
		if errorResponse.Error.Message != "" {
			return "", fmt.Errorf("AI Provider error %d: %s", resp.StatusCode, errorResponse.Error.Message)
		}
		return "", fmt.Errorf("AI Provider returned error status: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		content := result.Choices[0].Message.Content
		fmt.Printf("AI successfully generated tags: [%s]\n", content)
		return content, nil
	}

	return "", fmt.Errorf("no response from AI")
}