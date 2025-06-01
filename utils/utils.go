package utils

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func ApiCall(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API call: %w", err)
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}
	return responseBytes, nil
}

func NormalizeCityName(city string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`) // Matches non-alphanumeric characters
	return strings.ToLower(re.ReplaceAllString(city, " "))
}
