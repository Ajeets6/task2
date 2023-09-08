package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/.netlify/functions/api", apiHandler)
	http.ListenAndServe(":9000", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		operationCode := generateOperationCode()
		response := map[string]interface{}{
			"is_success":     true,
			"user_id":        generateUserID(),
			"operation_code": operationCode,
		}
		jsonResponse(w, response)
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err != nil {
			jsonResponse(w, map[string]interface{}{
				"is_success": false,
				"user_id":    generateUserID(),
				"message":    "Invalid JSON data",
			})
			return
		}

		inputData, ok := data["data"].([]interface{})
		if !ok {
			jsonResponse(w, map[string]interface{}{
				"is_success": false,
				"user_id":    generateUserID(),
				"message":    "Invalid input data",
			})
			return
		}

		var alphabets []string
		for _, char := range inputData {
			if s, isString := char.(string); isString {
				if len(s) == 1 && strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
					alphabets = append(alphabets, s)
				}
			}
		}

		var highestAlphabet string
		if len(alphabets) > 0 {
			highestAlphabet = findHighestAlphabet(alphabets)
		}

		response := map[string]interface{}{
			"is_success":     true,
			"user_id":        generateUserID(),
			"email":          "rk5532@srmist.edu.in",
			"roll_number":    "RA2011030020043",
			"numbers":        extractNumbers(inputData),
			"alphabets":      alphabets,
			"highest_alphabet": []string{highestAlphabet},
		}
		jsonResponse(w, response)
	default:
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
	}
}

func generateUserID() string {
	return "rohit_kumar_27062002"
}

func generateOperationCode() string {
	return "1"
}

func findHighestAlphabet(alphabets []string) string {
	highest := "A"
	for _, char := range alphabets {
		if char > highest {
			highest = char
		}
	}
	return highest
}

func extractNumbers(data []interface{}) []string {
	var numbers []string
	for _, item := range data {
		if s, isString := item.(string); isString {
			if _, err := json.Number(s).Float64(); err == nil {
				numbers = append(numbers, s)
			}
		}
	}
	return numbers
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
