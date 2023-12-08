package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order        []*Order
	Transactions []*Transaction
	OrderIn      chan *Order
	OrderOut     chan *Order
	WaitGroup    *sync.WaitGroup
}

func NewBook(orderIn, OrderOut chan *Order, waitGroup *sync.WaitGroup) *Book {
	return &Book{
		Order:        []*Order{},
		Transactions: []*Transaction{},
		OrderIn:      orderIn,
		OrderOut:     OrderOut,
		WaitGroup:    waitGroup,
	}
}

func (book *Book) Trade() {
	buyOrders := make(map[string]*Order_Queue)
	sellOrders := make(map[string]*Order_Queue)

	for order := range book.OrderIn {
		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		if order.OrderType == "BUY" {
			buyOrders[asset].Push(order)

			assetExists := sellOrders[asset].Len() > 0
			orderPriceIsEqualOrLower := order.Price >= sellOrders[asset].Order[0].Price

			if assetExists && orderPriceIsEqualOrLower {
				sellOrder := sellOrders[asset].Pop().(*Order)

				hasPeddingShares := sellOrder.PendingShares > 0

				if hasPeddingShares {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					book.AddTransaction(transaction, book.WaitGroup)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)

					book.OrderOut <- sellOrder
					book.OrderOut <- order

					if sellOrder.PendingShares > 0 {
						sellOrders[asset].Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders[asset].Push(order)

			assetExists := buyOrders[asset].Len() > 0
			orderPriceIsEqualOrHigher := order.Price <= buyOrders[asset].Order[0].Price

			if assetExists && orderPriceIsEqualOrHigher {
				buyOrder := buyOrders[asset].Pop().(*Order)

				hasPeddingShares := buyOrder.PendingShares > 0

				if hasPeddingShares {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					book.AddTransaction(transaction, book.WaitGroup)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)

					book.OrderOut <- buyOrder
					book.OrderOut <- order

					if buyOrder.PendingShares > 0 {
						buyOrders[asset].Push(buyOrder)
					}
				}
			}
		}
	}
}

func (book *Book) AddTransaction(transaction *Transaction, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares

	if buyingShares < sellingShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.SellOrderAddPeddingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.BuyOrderAddPeddingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseOrders()

	book.Transactions = append(book.Transactions, transaction)
}
