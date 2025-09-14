package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"playwithutf/services"
	"strconv"
	"strings"
)

type Operation string

const (
	EncodeOperation Operation = "encode"
	DecodeOperation Operation = "decode"
)

type ResponseBody struct {
	Success bool          `json:"success"`
	Code    int           `json:"code"`
	Data    *ResponseData `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type ResponseData struct {
	Operation Operation   `json:"operation,omitempty"`
	Input     string      `json:"input,omitempty"`
	Output    interface{} `json:"output,omitempty"`
}

type RequestBody struct {
	Operation Operation `json:"operation"`
	Input     string    `json:"input"`
}

func PlayWithUTF8Handler(w http.ResponseWriter, r *http.Request) {
	var request RequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, false, &ResponseData{}, err.Error())
		return
	}

	var responseData *ResponseData
	var handlerErr error

	switch request.Operation {
	case "encode":
		responseData, handlerErr = handleEncodeOperation(request.Input)
	case "decode":
		responseData, handlerErr = handleDecodeOperation(request.Input)
	default:
		writeJSONResponse(w, http.StatusBadRequest, false, &ResponseData{}, "Invalid operation. Use 'encode' or 'decode'.")
		return
	}

	if handlerErr != nil {
		writeJSONResponse(w, http.StatusBadRequest, false, &ResponseData{}, handlerErr.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, true, responseData, "")
}

func handleEncodeOperation(input string) (*ResponseData, error) {
	var allEncodedBytes []byte
	for _, r := range input {
		encoded := services.Utf8Encode(r)
		allEncodedBytes = append(allEncodedBytes, encoded...)
	}
	return &ResponseData{
		Operation: EncodeOperation,
		Input:     input,
		Output:    allEncodedBytes,
	}, nil
}

func handleDecodeOperation(input string) (*ResponseData, error) {
	parts := strings.Split(input, ",")
	var inputBytes []byte
	for _, part := range parts {
		hexValue, err := strconv.ParseUint(strings.TrimSpace(part), 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid hexadecimal value: %s", part)
		}
		inputBytes = append(inputBytes, byte(hexValue))
	}

	var decodedRunes []rune
	for i := 0; i < len(inputBytes); {
		r, size := services.Utf8Decode(inputBytes[i:])
		if size == 0 {
			return nil, fmt.Errorf("invalid UTF-8 sequence at position %d", i)
		}
		decodedRunes = append(decodedRunes, r)
		i += size
	}

	return &ResponseData{
		Operation: DecodeOperation,
		Input:     input,
		Output:    string(decodedRunes),
	}, nil
}

// writeJSONResponse is a helper function to standardize JSON responses.
func writeJSONResponse(w http.ResponseWriter, code int, success bool, data *ResponseData, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := ResponseBody{
		Success: success,
		Code:    code,
		Data:    data,
		Error:   error,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}
