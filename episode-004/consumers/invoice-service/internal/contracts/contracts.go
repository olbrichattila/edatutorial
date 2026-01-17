package contracts

// Episode 004 orderUUID added
type InvoiceRepository interface {
	CreateInvoice(orderUUID string) (int64, error)
}
