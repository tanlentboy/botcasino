package remote

import (
	"strconv"
	"sync"

	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/botcasino/protobuf/walletoserver"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// WalletServer 钱包服务
var WalletServer *WalletServerWrap

// WalletServerWrap 钱包服务
type WalletServerWrap struct {
	client WalletServerClient
}

var walletServerOnce sync.Once

// NewWalletServerForOnce 创建钱包服务连接
func NewWalletServerForOnce(address string, port int) {
	walletServerOnce.Do(func() {
		address = address + ":" + strconv.Itoa(port)
		conn, err := grpc.Dial(address,
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithTimeout(TimeOutDuration))
		if err != nil {
			logger.Panic(err)
		}

		WalletServer = &WalletServerWrap{
			client: WalletServerClient{
				WalletClient: walletoserver.NewWalletClient(conn),
			}}
	})
}

// RPC 获取RPC
func (wrap *WalletServerWrap) RPC() WalletServerClient {
	return wrap.client
}

// WalletServerClient 钱包服务客户端
type WalletServerClient struct {
	walletoserver.WalletClient
}

// Context 获取上下文
func (cli *WalletServerClient) Context(parent context.Context) context.Context {
	ctx, err := context.WithTimeout(parent, TimeOutDuration)
	if err != nil {
		return context.TODO()
	}
	return ctx
}
