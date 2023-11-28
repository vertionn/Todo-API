/*

	 main.go - This is the main file that holds the logic for the Todo API.

	 It defines the TodoStruct type representing individual todos, and the main router using the
	 go-chi framework. Endpoints for retrieving, creating, updating, and completing todos are implemented.
	 The application uses an in-memory slice 'Todos' to store todo items.

	 API Endpoints:
	 - GET /todos: Retrieve all todos or a message if none exist.
	 - POST /create/todo: Create a new todo.
	 - PUT /update/todo/{ID}: Update a todo by ID.
	 - PATCH /complete/{ID}: Mark a todo as complete by ID.
	 - DELETE "/delete/{ID}: Delete a todo a ID

	Author: Nathan
	Date: 24/11/23

*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type TodoStruct struct {
	ID          int
	Title       string `json:"title"`
	Description string `json:"description"`
	Complete    bool
}

func main() {
	var Todos []TodoStruct

	// r is a new chi router that will handle the HTTP routes.
	r := chi.NewRouter()

	// Use the Chi middleware Logger to log details about each incoming HTTP request.
	// This middleware captures and logs the request information, such as HTTP method, path, and duration.
	// It helps in debugging and monitoring the server's behavior without breaking the server itself.
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch len(Todos) {
		case 0:
			// If there are no Todos, return a message
			err := ReturnJSON(w, http.StatusOK, JsonResponse{
				Success: true,
				Message: "You have no todos. Try adding one.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		default:
			// If there are Todos, return them as JSON
			err := ReturnJSON(w, http.StatusOK, JsonResponse{
				Success: true,
				Todos:   Todos,
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	})

	r.Post("/create/todo", func(w http.ResponseWriter, r *http.Request) {

		var body TodoStruct

		// Create a new JSON decoder for the HTTP request body and disallow unknown
		// fields.
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		// Decode the HTTP request body into the 'body' struct.
		err := dec.Decode(&body)

		// Check for any decoding errors.
		if err != nil {
			// If there's an error, return a Bad Request response with an error message.
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "Invalid request data. Please ensure your request is properly formatted.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// add the new todo to the slice
		Todos = append(Todos, TodoStruct{ID: len(Todos) + 1, Title: body.Title, Description: body.Description, Complete: body.Complete})

		err = ReturnJSON(w, http.StatusBadRequest, JsonResponse{
			Success: true,
			Message: "Todo was created successfully.",
		})

		// If there is an error with encoding and sending the JSON back to the client, return the error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	})

	r.Put("/update/todo/{ID}", func(w http.ResponseWriter, r *http.Request) {

		var body TodoStruct

		// Create a new JSON decoder for the HTTP request body and disallow unknown
		// fields.
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		// Decode the HTTP request body into the 'body' struct.
		err := dec.Decode(&body)

		// Check for any decoding errors.
		if err != nil {
			// If there's an error, return a Bad Request response with an error message.
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "Invalid request data. Please ensure your request is properly formatted.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// grab the id from the url path and then convert it to an int and handle any errors
		ID := chi.URLParam(r, "ID")
		IDint, err := strconv.Atoi(ID)
		if err != nil {
			// If there's an error, return an error message.
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "There was a problem with the todo id, please fix it then try again.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var found bool

		// Iterate over the Todos to find the one with the matching ID.
		for i, k := range Todos {
			if k.ID == IDint {
				found = true

				// Update the title if it's different from the existing value and if the title in the body is non-empty.
				if k.Title != body.Title && body.Title != "" {
					Todos[i].Title = body.Title
				}

				// Update the description if it's different from the existing value or if it's explicitly an empty string.
				if k.Description != body.Description || body.Description == "" {
					Todos[i].Description = body.Description
				}

				// Respond with a success status code indicating that the todo was updated.
				w.WriteHeader(http.StatusNoContent)
				break
			}
		}

		// Check if a matching ID was not found.
		if !found {
			// return some json with the message telling the user we couldn't find the todo
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "We could not find any todo with this id, double check and try again.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	})

	r.Patch("/complete/{ID}", func(w http.ResponseWriter, r *http.Request) {

		// grab the id from the url path and then convert it to an int and handle any errors
		ID := chi.URLParam(r, "ID")
		IDint, err := strconv.Atoi(ID)
		if err != nil {
			// If there's an error, return an error message.
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "There was a problem with the todo id, please fix it then try again.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var found bool

		// Iterate over the Todos to find the one with the matching ID.
		for i, k := range Todos {
			if k.ID == IDint {
				found = true

				Todos[i].Complete = true

				w.WriteHeader(http.StatusNoContent)
				break
			}
		}

		// Check if a matching ID was not found.
		if !found {
			// return some json with the message telling the user we couldn't find the todo
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "We could not find any todo with this id, double check and try again.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

	})

	r.Delete("/delete/{ID}", func(w http.ResponseWriter, r *http.Request) {
		// Extract the ID from the URL parameters.
		ID := chi.URLParam(r, "ID")

		// Convert the ID to an integer.
		IDint, err := strconv.Atoi(ID)
		if err != nil {
			// If there's an error, return an error message.
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "There was a problem with the todo ID, please fix it and try again.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error.
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var found bool

		// Iterate over the Todos to find the one with the matching ID.
		for i, k := range Todos {
			if k.ID == IDint {
				found = true

				// Remove the todo from the Todos slice.
				Todos = append(Todos[:i], Todos[i+1:]...)

				// Respond with a success status code indicating that the todo was deleted.
				w.WriteHeader(http.StatusNoContent)
				break
			}
		}

		// Check if a matching ID was not found.
		if !found {
			// Return JSON with the message telling the user that no matching todo was found.
			err := ReturnJSON(w, http.StatusBadRequest, JsonResponse{
				Success:      false,
				ErrorMessage: "We could not find any todo with this ID, double-check and try again.",
			})

			// If there is an error with encoding and sending the JSON back to the client, return the error.
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	// Create a new HTTP server with the provided router
	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	// Set up a channel to listen for interrupts and gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		fmt.Println("Server is running on :3000")
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			// Only print the error if it's not a graceful shutdown error
			fmt.Println("Error:", err)
		}
	}()

	// Wait for an interrupt signal
	<-stop

	fmt.Println("Shutting down server...")

	// Create a context with a timeout to allow in-flight requests to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Error during server shutdown:", err)
	}

	fmt.Println("Server gracefully stopped.")
}
