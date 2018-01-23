package main

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/zhangpanyi/botcasino/config"
	"github.com/zhangpanyi/botcasino/envelopes"
	envelopescaches "github.com/zhangpanyi/botcasino/envelopes/caches"
	"github.com/zhangpanyi/botcasino/envelopes/expiretimer"
	"github.com/zhangpanyi/botcasino/envelopes/feessync"
	"github.com/zhangpanyi/botcasino/envelopes/groupchat"
	"github.com/zhangpanyi/botcasino/models"
	"github.com/zhangpanyi/botcasino/pusher"
	"github.com/zhangpanyi/botcasino/remote"
	"github.com/zhangpanyi/botcasino/service"
	"github.com/zhangpanyi/botcasino/storage"
	"github.com/zhangpanyi/botcasino/webrpc"
	"github.com/zhangpanyi/botcasino/withdraw"

	"github.com/vrecan/death"
	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/basebot/telegram/updater"
	"upper.io/db.v3/mysql"
)

// 读取红包封面
func readRedEnvelopesCover(files []string) error {
	cover := make([]groupchat.Cover, 0, len(files))
	for _, filename := range files {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		cover = append(cover, groupchat.Cover{
			FileData: data,
			FileName: filepath.Base(filename),
		})
	}
	groupchat.SetRedEnvelopesCoverForOnce(cover)
	return nil
}

func main() {
	// 加载配置文件
	config.LoadConfig("master.yml")

	// 初始化日志库
	serveCfg := config.GetServe()
	logger.CreateLoggerOnce(logger.DebugLevel, logger.InfoLevel)

	// 连接到数据库
	err := storage.Connect(serveCfg.BolTDBPath)
	if err != nil {
		logger.Panic(err)
	}

	// 连接MySQL数据库
	dbcfg := serveCfg.MySQL
	settings := mysql.ConnectionURL{
		Database: dbcfg.Database,
		Host:     dbcfg.Host,
		User:     dbcfg.User,
		Password: dbcfg.Password,
		Options:  dbcfg.Options,
	}
	err = models.Connect(settings, dbcfg.Conns)
	if err != nil {
		logger.Panic(err)
	}

	// 创建更新器
	botUpdater, err := updater.NewUpdater(serveCfg.Domain, serveCfg.APIWebsite)
	if err != nil {
		logger.Panic(err)
	}
	webrpc.InitRoute(botUpdater.GetRouter())

	// 读取红包封面图
	err = readRedEnvelopesCover(serveCfg.RedEnvelopesCover)
	if err != nil {
		logger.Panic(err)
	}

	// 连接钱包服务
	remote.NewWalletServerForOnce(serveCfg.WalletService.Address,
		serveCfg.WalletService.Port)

	// 同步转账手续费
	feessync.CheckFeeStatusAsync()

	// 运行转账服务
	withdraw.RunWithdrawServiceForOnce(6)

	// 启动红包机器人
	envelopescaches.CreateManagerForOnce(serveCfg.BucketNum)
	bot, err := botUpdater.AddHandler(serveCfg.Token, envelopes.NewUpdate)
	if err != nil {
		logger.Panic(err)
	}
	logger.Infof("Red envelope bot_id is: %d", bot.ID)

	pool, err := updater.NewPool(2048)
	if err != nil {
		logger.Panic(err)
	}
	expiretimer.StartTimerForOnce(bot, pool)

	// 创建消息推送器
	pusher.CreatePusherForOnce(pool)

	// 启动RPC服务
	port := strconv.Itoa(int(serveCfg.Port))
	address := serveCfg.BindAddress + ":" + port
	go service.RunService(address)

	// 启动更新服务器
	logger.Infof("Casino server started, grpc listen: %s:%d", serveCfg.BindAddress, serveCfg.Port)
	go func() {
		err = botUpdater.ListenAndServe(":443")
		if err != nil {
			logger.Panicf("Casino server failed to listen: %v", err)
		}
	}()

	// 捕捉退出信号
	d := death.NewDeath(syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM)
	d.WaitForDeathWithFunc(func() {
		storage.Close()
		logger.Info("Casino server stoped")
	})
}
