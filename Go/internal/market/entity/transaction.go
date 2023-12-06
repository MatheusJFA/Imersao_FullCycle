package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price

	return &Transaction{
		ID:           uuid.NewString(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func (transaction *Transaction) CalculateTotal(shares int, price float64) {
	transaction.Total = float64(shares) * price
}

func (transaction *Transaction) CloseBuyOrders() {
	if transaction.BuyingOrder.PendingShares == 0 {
		transaction.BuyingOrder.Status = "CLOSED"
	}
}

func (transaction *Transaction) CloseSellOrders() {
	if transaction.SellingOrder.PendingShares == 0 {
		transaction.SellingOrder.Status = "CLOSED"
	}
}

func (transaction *Transaction) CloseOrders() {
	transaction.CloseBuyOrders()
	transaction.CloseSellOrders()
}

func (transaction *Transaction) BuyOrderAddPeddingShares(shares int) {
	transaction.BuyingOrder.PendingShares += shares
}

func (transaction *Transaction) SellOrderAddPeddingShares(shares int) {
	transaction.SellingOrder.PendingShares += shares
}
