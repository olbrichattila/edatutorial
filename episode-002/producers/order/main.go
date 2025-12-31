// locally test with POST: http://localhost:8080/order
// Payload is in example-payload file
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/olbrichattila/edatutorial/shared/actions"
	event "github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
)

const topic = "order"

// DTO for orders to validate input
type order struct {
	UserID string  `json:"userId"`
	Email  string  `json:"email"`
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

		ord, err := translatePayloadToDTO(r)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		// Validate user data
		if err := validate(ord); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		actionPayload, err := translateToAction(ord)
		if err != nil {
			http.Error(w, "whoops something went wrong", http.StatusInternalServerError)
			return
		}

		// Publish event
		em.Publish(topic, actionPayload)

		// Return with accepted, no content 202
		w.WriteHeader(http.StatusAccepted)
	}
}

func validate(ord *order) error {
	// Homework: could validate user ID format and length
	if strings.TrimSpace(ord.UserID) == "" {
		return fmt.Errorf("user id required")
	}

	// Homework: validate email format as well
	if strings.TrimSpace(ord.Email) == "" {
		return fmt.Errorf("email required")
	}

	if len(ord.Items) == 0 {
		return fmt.Errorf("no items in order")
	}

	// Homework: could validate if all items have product number and quantity is larger then 0 and product id is not repeated
	return nil
}

func translatePayloadToDTO(r *http.Request) (*order, error) {
	var ord order
	err := json.NewDecoder(r.Body).Decode(&ord)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return &ord, nil
}

func translateToAction(ord *order) ([]byte, error) {
	// Create action from the validated input DTO
	ordItems := make([]actions.OrderItem, len(ord.Items))
	for i, it := range ord.Items {
		ordItems[i].ProductID = it.ProductID
		ordItems[i].Quantity = it.Quantity
	}

	envelope := actions.New[actions.OrderSentAction](actions.OrderSentAction{
		UserID: ord.UserID,
		Email:  ord.Email,
		Items:  ordItems,
	})

	// Create action payload from orderSent envelope
	return envelope.ToJSON()
}
