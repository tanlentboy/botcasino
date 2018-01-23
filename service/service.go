package service

import (
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/botcasino/models"
	"github.com/zhangpanyi/botcasino/protobuf/casinoserver"
	"github.com/zhangpanyi/botcasino/storage"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// 匹配用户ID
var reMathUserID *regexp.Regexp

func init() {
	var err error
	reMathUserID, err = regexp.Compile("^ *(\\d+) *$")
	if err != nil {
		logger.Panic(err)
	}
}

type casinoServer struct{}

// RunService 运行服务
func RunService(address string) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		logger.Panic(err)
	}

	svr := grpc.NewServer()
	casinoserver.RegisterCasinoServer(svr, new(casinoServer))
	svr.Serve(listen)
}

// SentNotice 转款通知
func (*casinoServer) SentNotice(ctx context.Context, req *casinoserver.SentNoticeRequest) (*casinoserver.SentNoticeReply, error) {
	go func() {
		logger.Warnf("On sent notice, Asset=%s, Amount=%d, BlockNum=%d, Memo=%s",
			req.GetAsset(), req.GetAmount(), req.GetBlockNum(), req.GetMemo())

		// 获取订单ID
		orderID, err := strconv.ParseInt(req.GetMemo(), 10, 64)
		if err != nil {
			logger.Warnf("On sent notice, invalid order id, %s", req.GetMemo())
			return
		}

		// 获取订单信息
		order, err := models.GetWithdrawOrder(orderID)
		if err != nil {
			logger.Warnf("On sent notice, nou found order id, %d", orderID)
			return
		}

		// 更新手续费
		models.UpdateWithdrawFee(orderID, uint32(req.GetFeeAmount()))

		// 插入操作记录
		desc := fmt.Sprintf("您提现*%.2f* *%s*已确认, 区块高度: *%d*", float64(req.GetAmount())/100.0,
			req.GetAsset(), req.GetBlockNum())
		models.InsertHistory(order.UserID, desc)
	}()
	return new(casinoserver.SentNoticeReply), nil
}

// ReceiveNotice 收款通知
func (*casinoServer) ReceiveNotice(ctx context.Context, req *casinoserver.ReceiveNoticeRequest) (*casinoserver.ReceiveNoticeReply, error) {
	go func() {
		logger.Warnf("On receive notice, Asset=%s, Amount=%d, BlockNum=%d, Memo=%s",
			req.GetAsset(), req.GetAmount(), req.GetBlockNum(), req.GetMemo())

		// 获取用户ID
		result := reMathUserID.FindStringSubmatch(req.GetMemo())
		if len(result) != 2 {
			logger.Warnf("On receive notice, not found user id, Asset=%s, Amount=%d, BlockNum=%d, Memo=%s",
				req.GetAsset(), req.GetAmount(), req.GetBlockNum(), req.GetMemo())
			return
		}

		userID, err := strconv.ParseInt(result[1], 10, 64)
		if err != nil {
			logger.Warnf("On receive notice, not found user id, Asset=%s, Amount=%d, BlockNum=%d, Memo=%s",
				req.GetAsset(), req.GetAmount(), req.GetBlockNum(), req.GetMemo())
			return
		}

		// 检查资产类型
		if req.GetAsset() != storage.BitCNYSymbol && req.GetAsset() != storage.BitUSDSymbol {
			logger.Warnf("On receive notice, nonsupport asset, Asset=%s, Amount=%d, BlockNum=%d, Memo=%s",
				req.GetAsset(), req.GetAmount(), req.GetBlockNum(), req.GetMemo())
			return
		}

		// 增加用户资产
		handler := storage.AssetStorage{}
		err = handler.Deposit(userID, req.GetAsset(), uint32(req.GetAmount()))
		if err != nil {
			logger.Errorf("On receive notice, deposit failure, UserID=%d, Asset=%s, Amount=%d, BlockNum=%d, Memo=%s, %v",
				userID, req.GetAsset(), req.GetAmount(), req.GetBlockNum(), req.GetMemo(), err)
			return
		}
		logger.Errorf("On receive notice, deposit success, UserID=%d, Asset=%s, Amount=%d, BlockNum=%d, Memo=%s",
			userID, req.GetAsset(), req.GetAmount(), req.GetBlockNum(), req.GetMemo())

		// 插入操作记录
		desc := fmt.Sprintf("您充值*%.2f* *%s*已确认, 区块高度: *%d*", float64(req.GetAmount())/100.0,
			req.GetAsset(), req.GetBlockNum())
		models.InsertHistory(userID, desc)
	}()
	return new(casinoserver.ReceiveNoticeReply), nil
}
