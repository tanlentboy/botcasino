package privatechat

import (
	"fmt"

	"github.com/zhangpanyi/botcasino/config"

	"github.com/zhangpanyi/basebot/history"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// Deposit 存款
type Deposit struct {
}

// Handle 消息处理
func (deposit *Deposit) Handle(bot *methods.BotExt, r *history.History, update *types.Update) {
	// 是否开放充值
	dynamicCfg := config.GetDynamic()
	fromID := update.CallbackQuery.From.ID
	if !dynamicCfg.AllowDeposit {
		reply := tr(fromID, "lng_priv_deposit_not_allow")
		bot.AnswerCallbackQuery(update.CallbackQuery, reply, false, "", 0)
	}

	// 回复充值步骤
	serveCfg := config.GetServe()
	reply := fmt.Sprintf(tr(fromID, "lng_priv_deposit_say"), serveCfg.Account, fromID)
	menus := [...]methods.InlineKeyboardButton{
		methods.InlineKeyboardButton{
			Text:         tr(fromID, "lng_back_superior"),
			CallbackData: "/main/",
		},
	}
	markup := methods.MakeInlineKeyboardMarkupAuto(menus[:], 1)
	bot.AnswerCallbackQuery(update.CallbackQuery, "", false, "", 0)
	bot.EditMessageReplyMarkup(update.CallbackQuery.Message, reply, true, markup)
}

// 消息路由
func (deposit *Deposit) route(bot *methods.BotExt, query *types.CallbackQuery) Handler {
	return nil
}
