package utils

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
)

type VirusTotalResponse struct {
	Data struct {
		Attributes struct {
			TotalVotes struct {
				Harmless int `json:"harmless"`
				Malicious int `json:"malicious"`
			} `json:"total_votes"`
		} `json:"attributes"`
	} `json:"data"`
}

func CheckVirusTotalHash(fileHash,virustotalAPIKey string) (string, error) {

	url := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", fileHash)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("x-apikey", virustotalAPIKey)
	req.Header.Add("accept", "application/json")

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

	var virustotalResponse VirusTotalResponse
	if err := json.Unmarshal(body, &virustotalResponse); err != nil {
		return "", err
	}

	if virustotalResponse.Data.Attributes.TotalVotes.Malicious > 0 {
		return "Malicious", nil
	}
	return "Harmless", nil
}
