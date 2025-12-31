package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/notification"
	invoiceContracts "producer.example/internal/contracts"
	"producer.example/internal/repositories/invoice"
)

const (
	topic    = "invoicecreated"
	consumer = "sendinvoiceemail"

	logTopic = "logmessagecreated"
)

func main() {
	eventManager := event.New()

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	invoiceRepository := invoice.New(db)

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		fmt.Println(log)
		evt.Publish(logTopic, []byte(log))

		invoiceCreatedAction, err := actions.FromJSON[actions.InvoiceCreatedAction](msg)

		// Get invoice head and body from db
		head, items, err := retrieveInvoiceDetails(invoiceRepository, invoiceCreatedAction.Payload.ID)
		if err != nil {
			evt.Publish(logTopic, []byte("cannot fetch invoice: "+err.Error()))
			return err
		}

		// Generate HTML to send
		html, total, err := invoiceAsHTML(head, items)
		if err != nil {
			evt.Publish(logTopic, []byte("cannot crate invoice html: "+err.Error()))
			return err
		}

		// Finalize email body
		emailBody := fmt.Sprintf(`<html><body><h2>Invoice</h2>%s</html><p>Total: <b>%.2f</b></p>`, html, total)

		// Send email
		err = notification.SendEmail(string(head["email"].([]byte)), "Your invoice", emailBody)
		if err != nil {
			evt.Publish(logTopic, []byte("send email error: "+err.Error()))
			return err
		}

		return nil
	})
}

func retrieveInvoiceDetails(invoiceRepository invoiceContracts.InvoiceRepository, invoiceId int64) (map[string]any, []map[string]any, error) {
	head, err := invoiceRepository.Head(invoiceId)
	if err != nil {
		return nil, nil, err
	}

	// Get invoice items from db
	items, err := invoiceRepository.Items(invoiceId)
	if err != nil {
		return nil, nil, err
	}

	return head, items, nil
}

func invoiceAsHTML(head map[string]any, items []map[string]any) (string, float64, error) {
	headAsHTML := headAsHTML(head)
	itemsAsHTML, total, err := itemsAsHTML(items)
	if err != nil {
		return "", 0, err
	}

	return headAsHTML + itemsAsHTML, total, nil
}

func headAsHTML(head map[string]any) string {
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
		head["id"],
		head["user_id"],
		head["email"],
	)
}

func itemsAsHTML(items []map[string]any) (string, float64, error) {
	itemsAsHTML := `<table style=" border-collapse: collapse; border: 1px solid black; width: 100%">`
	total := 0.0
	for _, item := range items {
		itemsAsHTML += `
			<tr>` +
			fmt.Sprintf(`<td style="border: 1px solid black; padding: 3px 6px;">%s</td>`, item["product_id"]) +
			fmt.Sprintf(`<td style="border: 1px solid black; padding: 3px 6px;">%d</td>`, item["quantity"]) +
			fmt.Sprintf(`<td style="border: 1px solid black; padding: 3px 6px;">%s</td>`, item["price"]) +
			`</tr>`

		price, err := strconv.ParseFloat(string(item["price"].([]byte)), 10)
		if err != nil {
			return "", 0, err
		}

		total += float64(item["quantity"].(int64)) * price
	}
	itemsAsHTML += "</table>"

	return itemsAsHTML, total, nil
}
