package transformer

import (
	"github.com/MatheusJFA/Imersao_FullCycle/Go/internal/market/dto"
	"github.com/MatheusJFA/Imersao_FullCycle/Go/internal/market/entity"
)

func TransformInput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)

	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}

	return order
}

func TransformOutput(order *entity.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderID:    order.ID,
		InvestorID: order.Investor.ID,
		AssetID:    order.Asset.ID,
		Shares:     order.Shares,
		Status:     order.Status,
		Partial:    order.PendingShares,
		OrderType:  order.OrderType,
	}

	var transactionsOutput []*dto.TransactionOutput

	for _, transaction := range order.Transactions {
		total := transaction.Shares - transaction.SellingOrder.PendingShares
		transactionOutput := &dto.TransactionOutput{
			TransactionID: transaction.ID,
			BuyerID:       transaction.BuyingOrder.ID,
			SellerID:      transaction.SellingOrder.ID,
			AssetID:       transaction.SellingOrder.Asset.ID,
			Price:         transaction.Price,
			Shares:        total,
		}

		transactionsOutput = append(transactionsOutput, transactionOutput)
	}

	output.TransactionsOutput = transactionsOutput

	return output

}
