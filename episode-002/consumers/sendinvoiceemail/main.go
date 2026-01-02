package main

import (
	"fmt"
	"html"
	"os"
	"strconv"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
	"github.com/olbrichattila/edatutorial/shared/notification"
	invoiceContracts "producer.example/internal/contracts"
	"producer.example/internal/repositories/invoice"
)

const (
	topic    = "invoicecreated"
	consumer = "sendinvoiceemail"
)

func main() {
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := eventlogger.New(eventManager)

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Printf("Error closing database: %v\n", closeErr)
		}
	}()

	invoiceRepository := invoice.New(db)

	err = eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.Info(log)

		invoiceCreatedAction, err := actions.FromJSON[actions.InvoiceCreatedAction](msg)
		if err != nil {
			logger.Error("cannot parse JSON: " + err.Error())
			return err
		}

		// Get invoice head and body from db
		head, items, err := retrieveInvoiceDetails(invoiceRepository, invoiceCreatedAction.Payload.ID)
		if err != nil {
			logger.Error("cannot fetch invoice: " + err.Error())
			return err
		}

		// Generate HTML to send
		html, total, err := invoiceAsHTML(head, items)
		if err != nil {
			logger.Error("cannot fetch invoice HTML: " + err.Error())
			return err
		}

		// Finalize email body
		emailBody := fmt.Sprintf(`<html><body><h2>Invoice</h2>%s</html><p>Total: <b>%.2f</b></p>`, html, total)

		// Send email
		err = notification.SendEmail(string(head["email"].([]byte)), "Your invoice", emailBody)
		if err != nil {
			logger.Error("send email error: " + err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error starting consumer: %v\n", err)
		os.Exit(1)
	}
}

func retrieveInvoiceDetails(invoiceRepository invoiceContracts.InvoiceRepository, invoiceID int64) (map[string]any, []map[string]any, error) {
	head, err := invoiceRepository.Head(invoiceID)
	if err != nil {
		return nil, nil, err
	}

	// Get invoice items from db
	items, err := invoiceRepository.Items(invoiceID)
	if err != nil {
		return nil, nil, err
	}

	return head, items, nil
}

func invoiceAsHTML(head map[string]any, items []map[string]any) (string, float64, error) {
	headHTML, err := headAsHTML(head)
	if err != nil {
		return "", 0, err
	}

	itemsHTML, total, err := itemsAsHTML(items)
	if err != nil {
		return "", 0, err
	}

	return headHTML + "\n" + itemsHTML, total, nil
}

func headAsHTML(head map[string]any) (string, error) {
	id, ok := head["id"].(int64)
	if !ok {
		return "", fmt.Errorf("id is not an int64")
	}

	userID, ok := head["user_id"].([]byte)
	if !ok {
		return "", fmt.Errorf("userID is not a string")
	}

	email, ok := head["email"].([]byte)
	if !ok {
		return "", fmt.Errorf("email is not a string")
	}

	return fmt.Sprintf(
		`<table style=" border-collapse: collapse; border: 1px solid black; width: 100%%">
			<tr>
				<td>Invoice ID:</td>
				<td colspan="2">%d</td>
			</tr>
			<tr>
				<td>User ID:</td>
				<td colspan="2">%s</td>
			</tr>
			<tr>
				<td>Email:</td>
				<td>%s</td>
			</tr>
			</table>`,
		id,
		html.EscapeString(string(userID)),
		html.EscapeString(string(email)),
	), nil
}

func itemsAsHTML(items []map[string]any) (string, float64, error) {
	itemsAsHTML := `<table style=" border-collapse: collapse; border: 1px solid black; width: 100%">`
	total := 0.0
	for _, item := range items {

		productID, ok := item["product_id"].([]byte)
		if !ok {
			return "", 0, fmt.Errorf("product_id is not a string")
		}

		quantity, ok := item["quantity"].(int64)
		if !ok {
			return "", 0, fmt.Errorf("quantity is not an int")
		}

		priceAsByteSlice, ok := item["price"].([]byte)
		if !ok {
			return "", 0, fmt.Errorf("price is not a string")
		}

		price, err := strconv.ParseFloat(string(priceAsByteSlice), 10)
		if err != nil {
			return "", 0, err
		}

		itemsAsHTML += `
			<tr>` +
			fmt.Sprintf(`<td style="border: 1px solid black; padding: 3px 6px;">%s</td>`, html.EscapeString(string(productID))) +
			fmt.Sprintf(`<td style="border: 1px solid black; padding: 3px 6px;">%d</td>`, quantity) +
			fmt.Sprintf(`<td style="border: 1px solid black; padding: 3px 6px;">%0.2f</td>`, price) +
			`</tr>`

		total += float64(item["quantity"].(int64)) * price
	}
	itemsAsHTML += "</table>"

	return itemsAsHTML, total, nil
}
