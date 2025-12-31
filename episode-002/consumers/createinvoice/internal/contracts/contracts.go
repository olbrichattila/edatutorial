package contracts

type InvoiceRepository interface {
	CreateInvoice(orderId int64) (int64, error)
}
