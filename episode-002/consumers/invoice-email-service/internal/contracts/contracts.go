package contracts

type InvoiceRepository interface {
	Head(invoiceID int64) (map[string]any, error)
	Items(invoiceID int64) ([]map[string]any, error)
}
