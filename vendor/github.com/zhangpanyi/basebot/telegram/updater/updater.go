package updater

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// 匹配webhook
var reMathWebhook *regexp.Regexp

func init() {
	var err error
	reMathWebhook, err = regexp.Compile("^/([a-zA-Z0-9:_-]+)/$")
	if err != nil {
		log.Panicln(err)
	}
}

// NewUpdater 创建更新器
func NewUpdater(domain string, apiwebsite string) (*Updater, error) {
	certificate, err := ioutil.ReadFile("certs/server.crt")
	if err != nil {
		return nil, err
	}

	updater := Updater{
		domain:      domain,
		apiwebsite:  apiwebsite,
		certificate: certificate,
		queue:       NewQueue(1024),
		handlers:    make(map[string]Handler),
		bots:        make(map[string]methods.BotExt),
	}
	pool, err := NewPool(2048)
	if err != nil {
		return nil, err
	}
	updater.pool = pool
	return &updater, nil
}

// Handler 处理函数
type Handler func(bot *methods.BotExt, update *types.Update)

// Updater 更新器
type Updater struct {
	domain        string                    // 服务域名
	apiwebsite    string                    // 机器人API服务网址
	certificate   []byte                    // 证书信息
	router        http.ServeMux             // 路由器
	bots          map[string]methods.BotExt // 机器人信息
	botMutex      sync.RWMutex              // 机器人信息锁
	handlers      map[string]Handler        // Token处理模块
	handlersMutex sync.RWMutex              // Token处理模块锁
	pool          *Pool                     // 工作队列
	queue         *Queue                    // 消息队列
}

// GetRouter 获取路由器
func (updater *Updater) GetRouter() *http.ServeMux {
	return &updater.router
}

// AddHandler 添加处理模块
func (updater *Updater) AddHandler(token string, handler Handler) (*methods.BotExt, error) {
	// 获取机器人
	bot, err := methods.GetMe(updater.apiwebsite, token)
	if err != nil {
		return nil, err
	}

	// 重设webhhok
	url := "https://" + updater.domain + "/" + token + "/"
	allowedUpdates := [...]string{"message", "callback_query"}
	err = bot.SetWebhook(url, updater.certificate, 40, allowedUpdates[:])
	if err != nil {
		return nil, err
	}

	// 插入机器人
	updater.botMutex.Lock()
	updater.bots[token] = *bot
	updater.botMutex.Unlock()

	// 插入处理模块
	updater.handlersMutex.Lock()
	updater.handlers[token] = handler
	updater.handlersMutex.Unlock()

	// 注册路由
	pattern := fmt.Sprintf("/%s/", token)
	updater.router.HandleFunc(pattern, updater.handleFunc)
	return bot, nil
}

// RemoveHandler 移除处理模块
func (updater *Updater) RemoveHandler(token string) {
	// 删除机器人
	updater.botMutex.Lock()
	delete(updater.bots, token)
	updater.botMutex.Unlock()

	// 删除处理模块
	updater.handlersMutex.Lock()
	delete(updater.handlers, token)
	updater.handlersMutex.Unlock()

	// 注销路由
	pattern := fmt.Sprintf("/%s/", token)
	updater.router.HandleFunc(pattern, nil)
}

// ListenAndServe 监听并服务
func (updater *Updater) ListenAndServe(addr string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: &updater.router,
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},
	}
	return s.ListenAndServeTLS("certs/server.crt", "certs/server.key")
}

// HTTP处理函数
func (updater *Updater) handleFunc(res http.ResponseWriter, req *http.Request) {
	// 获取token
	submatch := reMathWebhook.FindStringSubmatch(req.URL.Path)
	if len(submatch) < 2 {
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("Telegram updater: invalid url path, %v", req.URL.Path)
		return
	}
	token := submatch[1]

	// 获取机器人
	var ok bool
	var bot methods.BotExt
	updater.botMutex.RLock()
	bot, ok = updater.bots[token]
	if !ok {
		updater.botMutex.RUnlock()
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Telegram updater: invalid token, %v", token)
		return
	}
	updater.botMutex.RUnlock()

	// 处理更新数据
	jsb, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("Telegram updater: read err, %v", err)
		return
	}

	var update types.Update
	err = json.Unmarshal(jsb, &update)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Printf("Telegram updater: read err, %v", err)
		return
	}

	// 获取处理程序
	var handler Handler
	updater.handlersMutex.RLock()
	handler, ok = updater.handlers[token]
	if !ok {
		updater.handlersMutex.RUnlock()
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Telegram updater: handler not found, %v", token)
		return
	}
	updater.handlersMutex.RUnlock()

	// 分发消息处理
	updater.queue.Put(future{
		handler: handler,
		bot:     &bot,
		update:  &update,
		pool:    updater.pool,
	})

	res.WriteHeader(http.StatusOK)
}
