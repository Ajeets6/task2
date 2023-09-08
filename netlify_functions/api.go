package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// APIResponse defines the response structure.
type APIResponse struct {
	IsSuccess       bool     `json:"is_success"`
	UserID          string   `json:"user_id"`
	OperationCode   string   `json:"operation_code"`
	Message         string   `json:"message,omitempty"`
	Email           string   `json:"email,omitempty"`
	RollNumber      string   `json:"roll_number,omitempty"`
	Numbers         []string `json:"numbers,omitempty"`
	Alphabets       []string `json:"alphabets,omitempty"`
	HighestAlphabet []string `json:"highest_alphabet,omitempty"`
}

func main() {
	lambda.Start(apiHandler)
}

func apiHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		operationCode := generateOperationCode()
		response := APIResponse{
			IsSuccess:     true,
			UserID:        generateUserID(),
			OperationCode: operationCode,
		}
		return createResponse(http.StatusOK, response), nil

	case "POST":
		var inputData struct {
			Data []string `json:"data"`
		}
		if err := json.Unmarshal([]byte(request.Body), &inputData); err != nil {
			response := APIResponse{
				IsSuccess: false,
				UserID:    generateUserID(),
				Message:   "Invalid JSON data",
			}
			return createResponse(http.StatusBadRequest, response), nil
		}

		var alphabets []string
		for _, char := range inputData.Data {
			if len(char) == 1 && strings.ContainsAny(char, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
				alphabets = append(alphabets, char)
			}
		}

		var highestAlphabet string
		if len(alphabets) > 0 {
			highestAlphabet = findHighestAlphabet(alphabets)
		}

		response := APIResponse{
			IsSuccess:     true,
			UserID:        generateUserID(),
			Email:         "rk5532@srmist.edu.in",
			RollNumber:    "RA2011030020043",
			Numbers:       extractNumbers(inputData.Data),
			Alphabets:     alphabets,
			HighestAlphabet: []string{highestAlphabet},
		}
		return createResponse(http.StatusOK, response), nil

	default:
		response := APIResponse{
			IsSuccess: false,
			UserID:    generateUserID(),
			Message:   "Invalid HTTP method",
		}
		return createResponse(http.StatusMethodNotAllowed, response), nil
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

func extractNumbers(data []string) []string {
	var numbers []string
	for _, item := range data {
		if _, err := json.Number(item).Float64(); err == nil {
			numbers = append(numbers, item)
		}
	}
	return numbers
}

func createResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	responseBody, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}
}
