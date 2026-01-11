package contracts

type InvoiceRepository interface {
	CreateInvoice(orderID int64) (int64, error)
}
