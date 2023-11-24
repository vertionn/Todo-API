/*
   json.go - JSON Handling Module

   This file provides JSON handling functionalities for encoding and sending
   JSON responses to clients in a standardized format. It utilizes the 'jsoniter'
   library, which is compatible with the standard library, for efficient JSON encoding.

   The primary components include:
   - NewJson: A jsoniter configuration compatible with the standard library.
   - JsonResponse: Represents the structure of JSON responses sent to the client.
   - ReturnJSON: Function to encode and send a JSON response with proper content type
     and HTTP status code. Handles errors during encoding or writing.

   Author: Nathan
   Date: 24/11/23
*/

package main

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

// NewJson is a jsoniter configuration compatible with the standard library.
var NewJson = jsoniter.ConfigCompatibleWithStandardLibrary

// JsonResponse represents the structure of JSON responses sent to the client.
type JsonResponse struct {

	// Success indicates if the request was successful.
	Success bool `json:"success"`

	// Message is an optional message for success.
	Message string `json:"message,omitempty"`

	// ErrorMessage is an optional error message for failure.
	ErrorMessage string `json:"error_message,omitempty"`

	// Todos is an optional todo array to show all todos.
	Todos []TodoStruct `json:"todos,omitempty"`
}

// ReturnJSON encodes and sends a JSON response using the 'JsonResponse' structure.
// It sets the appropriate Content-Type and HTTP status code before encoding and sending the JSON.
// If there are any errors during the encoding or writing process, an error is returned.
func ReturnJSON(w http.ResponseWriter, statusCode int, response JsonResponse) error {
	// Set the Content Type
	w.Header().Set("Content-Type", "application/json")
	// Set the status code
	w.WriteHeader(statusCode)

	// Encode the JSON response and send it to the ResponseWriter
	err := NewJson.NewEncoder(w).Encode(response)
	if err != nil {
		// If an error occurs during encoding, return the error for further handling.
		return err
	}

	// No errors occurred during encoding and writing, return nil.
	return nil
}
