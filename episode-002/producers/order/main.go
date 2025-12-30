// locally test with POST: http://localhost:8080/order
// Payload is in example-payload file
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	event "github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
)

const topic = "order"

type order struct {
	UserID string  `json:"userId"`
	Items  []items `json:"items"`
}

type items struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func main() {
	eventManager := event.New()

	http.HandleFunc("/order", orderHandler(eventManager))

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func orderHandler(em contracts.EventManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If not POST error
		if r.Method != http.MethodPost {
			http.Error(w, "not a POST request", http.StatusBadRequest)
			return
		}

		// Marshal JSON payload
		var ord order
		err := json.NewDecoder(r.Body).Decode(&ord)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate user data
		if err := validate(&ord); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Re-marshal to event payload, here for simplicity I use the same JSON
		eventPayload, err := json.Marshal(ord)
		if err != nil {
			http.Error(w, "whoops something went wrong", http.StatusInternalServerError)
			return
		}

		// Publish event
		em.Publish(topic, eventPayload)

		// Return with accepted, no content 202
		w.WriteHeader(http.StatusAccepted)
	}
}

func validate(ord *order) error {
	// Homework: could validate user ID format and length
	if strings.TrimSpace(ord.UserID) == "" {
		return fmt.Errorf("user id required")
	}

	if len(ord.Items) == 0 {
		return fmt.Errorf("no items in order")
	}

	// Homework: could validate if all items have product number and quantity is larger then 0 and product id is not repeated

	return nil
}
