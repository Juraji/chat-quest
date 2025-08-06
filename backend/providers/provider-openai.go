package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type openAIProvider struct {
	connectionProfile *ConnectionProfile
}

type openAIListResponse[T interface{}] struct {
	Data []T `json:"data"`
}

type modelData struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

type embeddingRequestData struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type embeddingData struct {
	Embedding Embeddings `json:"embedding"`
}

func (o *openAIProvider) doRequest(method, path string, body io.Reader, v interface{}) error {
	url := o.connectionProfile.BaseUrl + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+o.connectionProfile.ApiKey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if v == nil {
		return nil
	} // caller only wants status
	return json.Unmarshal(data, v)
}

func (o *openAIProvider) getAvailableModels() ([]*LlmModel, error) {
	var resp openAIListResponse[modelData]
	if err := o.doRequest("GET", "/models", nil, &resp); err != nil {
		return nil, fmt.Errorf("get models: %w", err)
	}

	models := make([]*LlmModel, len(resp.Data))
	for i, d := range resp.Data {
		models[i] = DefaultLlmModel(o.connectionProfile.ID, d.ID)
	}
	return models, nil
}

func (o *openAIProvider) generateEmbeddings(input, modelID string) (Embeddings, error) {
	reqBody := embeddingRequestData{Input: input, Model: modelID}
	bodyBytes, _ := json.Marshal(reqBody)

	var resp openAIListResponse[embeddingData]
	if err := o.doRequest("POST", "/embeddings", bytes.NewReader(bodyBytes), &resp); err != nil {
		return nil, fmt.Errorf("generate embeddings: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, errors.New("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}
