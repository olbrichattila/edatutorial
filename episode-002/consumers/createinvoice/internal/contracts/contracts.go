package contracts

type InvoiceRepository interface {
	CreateInvoice(orderId int64) error
}
