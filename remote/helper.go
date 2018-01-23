package remote

import (
	"context"
	"errors"
	"strconv"

	"github.com/zhangpanyi/botcasino/protobuf/walletoserver"
)

// GetFees 获取手续费
func GetFees(assets []string) ([]uint32, error) {
	request := walletoserver.GetFeeRequest{
		Assets: assets,
	}
	rpc := WalletServer.RPC()
	result, err := rpc.GetFees(rpc.Context(context.TODO()), &request)
	if err != nil {
		return nil, err
	}

	if !result.GetOk() {
		return nil, errors.New(result.GetReason())
	}

	fees := make([]uint32, 0, len(result.GetFees()))
	for _, fee := range result.GetFees() {
		fees = append(fees, fee.GetFee())
	}
	return fees, nil
}

// Transfer 转账操作
func Transfer(orderID int64, to, assetID string, amount uint32) error {
	memo := strconv.FormatInt(orderID, 10)
	request := walletoserver.TransferRequest{
		To:      to,
		AssetId: assetID,
		Amount:  amount,
		Memo:    memo,
	}
	rpc := WalletServer.RPC()
	result, err := rpc.Transfer(rpc.Context(context.TODO()), &request)
	if err != nil {
		return err
	}

	if !result.GetOk() {
		return errors.New(result.GetReason())
	}
	return nil
}
