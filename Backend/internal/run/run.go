package run

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nikhilsiwach28/MyCode.git/models"
	"github.com/nikhilsiwach28/MyCode.git/queue"
	"github.com/nikhilsiwach28/MyCode.git/redis"
	"gopkg.in/Shopify/sarama.v1"
)

// HandleRun handles the execution of a submission request.
func HandleRun(ws *websocket.Conn, queueService queue.QueueService, redisService *redis.RedisService) {
	defer ws.Close()

	// Read message from WebSocket
	_, message, err := ws.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	// Unmarshal the JSON message into a submission request
	var request models.CreateSubmissionAPIRequest
	if err := json.Unmarshal(message, &request); err != nil {
		log.Println("Error parsing JSON:", err)
		return
	}

	// Convert request to submission
	submission := request.ToSubmissions()

	// Marshal submission to JSON
	submissionJSON, err := json.Marshal(submission)
	if err != nil {
		log.Println("Error marshaling submission to JSON:", err)
		return
	}

	// Send submission to the execution queue
	if err := queueService.SendMessage("toBeExecuted", string(submissionJSON)); err != nil {
		log.Println("Error sending message to queue:", err)
		return
	}

	// Store submission ID and input file in Redis
	if err := redisService.Set(submission.ID.String(), string(request.InputFile)); err != nil {
		log.Println("Error setting value in Redis:", err)
		return
	}

	// Listen for execution results
	handleExecutionResults(ws, queueService, redisService)
}

// handleExecutionResults listens for execution results and sends them to the WebSocket.
func handleExecutionResults(ws *websocket.Conn, queueService queue.QueueService, redisService *redis.RedisService) {
	// Subscribe to the executedCode queue
	messages, errors := queueService.Subscribe("executedCode")

	// Handle messages and errors
	for {
		select {
		case msg := <-messages:
			// Handle execution result message
			handleExecutionResultMessage(ws, msg, redisService)
			return
		case err := <-errors:
			// Handle queue error
			fmt.Println("Error in queue:", err)
			return
		}
	}
}

// handleExecutionResultMessage handles an execution result message.
func handleExecutionResultMessage(ws *websocket.Conn, msg *sarama.ConsumerMessage, redisService *redis.RedisService) {
	// Check message format
	if msg.Key == nil || msg.Value == nil {
		fmt.Println("Invalid message format")
		return
	}

	// Retrieve output from Redis using the output key
	output, err := redisService.Get(string(msg.Value))
	if err != nil {
		fmt.Println("Error retrieving output from Redis:", err)
		return
	}

	// Write output to WebSocket
	if err := ws.WriteMessage(websocket.TextMessage, []byte(output)); err != nil {
		log.Println("Error writing message to WebSocket:", err)
		return
	}

	// Close WebSocket after sending the response
	return
}
