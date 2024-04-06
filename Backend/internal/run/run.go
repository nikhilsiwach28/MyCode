package run

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nikhilsiwach28/MyCode.git/models"
	"github.com/nikhilsiwach28/MyCode.git/queue"
	"github.com/nikhilsiwach28/MyCode.git/redis"
)

func HandleRun(ws *websocket.Conn, queueService queue.QueueService, redisService *redis.RedisService) {
	defer ws.Close()

	_, message, err := ws.ReadMessage()
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

	var request models.CreateSubmissionAPIRequest
	if err := json.Unmarshal(message, &request); err != nil {
		log.Println("Error parsing JSON:", err)
		return
	}

	inputFile := request.InputFile
	submission := request.ToSubmissions()

	fmt.Printf("Received message type %d: %s\n", inputFile, submission)

	submissionJSON, err := json.Marshal(submission)
	if err != nil {
		log.Println("Error marshaling submission to JSON:", err)
		return
	}

	err = queueService.SendMessage("toBeExecuted", string(submissionJSON))
	if err != nil {
		log.Println("Error sending message to queue:", err)
		return
	}
	// Store submission ID and input file in Redis
	err = redisService.Set(submission.ID.String(), string(inputFile))
	if err != nil {
		log.Println("Error setting value in Redis:", err)
		return
	}

	// Subscribe to the executedCode queue
	messages, errors := queueService.Subscribe("executedCode")

	// Handle the messages and errors
	for {
		select {
		case msg := <-messages:
			// Check if the message contains the required fields
			if msg.Key == nil || msg.Value == nil {
				fmt.Println("Invalid message format")
				continue
			}

			// Retrieve the output from Redis using the output key
			output, err := redisService.Get(string(msg.Value))
			if err != nil {
				fmt.Println("Error retrieving output from Redis:", err)
				continue
			}

			err = ws.WriteMessage(websocket.TextMessage, []byte(output))
			if err != nil {
				log.Println("Error writing message to WebSocket:", err)
				return
			}
			// Close WebSocket after sending the response
			return
		case err := <-errors:
			fmt.Println("Error in queue:", err)
			// Handle the error from the queue
			return
		}
	}
}
