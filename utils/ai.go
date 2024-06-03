package utils

import (
	"autojob/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func askGPT(message string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	systemPrompt := `You are a job matching assistant.
	Your task is to evaluate whether the user cover letter and user job history match the requirements of a job description.
	You will receive a JSON payload with a user cover letter, user job history, job title, and job description.
	Based on this information, you will determine if there is a match (isMatch: true or false).
	If isMatch is true, you will also provide a custom cover letter tailored to the job description.

	The response should be in JSON format with the following structure:
	{
		"isMatch": boolean,
		"coverLetter": string
	}

	If isMatch is false, the coverLetter should be an empty string.

	Consider the following when making your decision:
	- Relevance of job history to the job description
	- Skills and experiences mentioned
	- Any other relevant information

	Make sure the cover letter is professional, concise, and highlights the candidate's strengths in relation to the job description.`

	requestBody := models.ChatCompletionRequest{
		Model: "gpt-4o",
		ResponseFormat: map[string]string{
			"type": "json_object",
		},
		Messages: []models.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: message},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var response models.ChatCompletionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
