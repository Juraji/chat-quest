package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type openAIProvider struct {
	ConnectionProfile *ConnectionProfile
}

type modelListResponse struct {
	Object string      `json:"object"`
	Data   []modelData `json:"data"`
}

type modelData struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

func (o *openAIProvider) getAvailableModels() ([]*LlmModel, error) {
	url := o.ConnectionProfile.BaseUrl + "/models"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+o.ConnectionProfile.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response modelListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	models := make([]*LlmModel, len(response.Data))
	for i, data := range response.Data {
		models[i] = DefaultLlmModel(o.ConnectionProfile.ID, data.ID)
	}

	return models, nil
}
