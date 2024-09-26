package IA

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Estructura para la respuesta de Cohere
type CohereResponse struct {
	Generations []struct {
		Text string `json:"text"`
	} `json:"generations"`
}

type CohereJSONResponse struct {
	Accion   string `json:"accion"`
	Producto string `json:"producto"`
	Cantidad int    `json:"cantidad"`
}

func GetCohereResponse(prompt, apiKey string) (CohereJSONResponse, error) {
	var response CohereJSONResponse
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":      "command-xlarge-nightly",
		"prompt":     prompt,
		"max_tokens": 100,
	})
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("POST", "https://api.cohere.ai/v1/generate", bytes.NewBuffer(requestBody))
	if err != nil {
		return response, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return response, fmt.Errorf("Error en la respuesta de Cohere: %s", string(body))
	}

	var cohereResponse CohereResponse
	err = json.NewDecoder(resp.Body).Decode(&cohereResponse)
	if err != nil {
		return response, err
	}

	if len(cohereResponse.Generations) > 0 {
		// Aquí intentamos decodificar el texto generado como JSON
		err := json.Unmarshal([]byte(cohereResponse.Generations[0].Text), &response)
		if err != nil {
			return response, fmt.Errorf("Error al parsear el JSON: %v", err)
		}
	}

	return response, nil
}
