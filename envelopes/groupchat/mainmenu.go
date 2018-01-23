package groupchat

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zhangpanyi/botcasino/config"
	"github.com/zhangpanyi/botcasino/models"
	"github.com/zhangpanyi/botcasino/storage"

	"github.com/zhangpanyi/basebot/history"
	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// Handler 消息处理
type Handler interface {
	route(*methods.BotExt, *types.CallbackQuery) Handler
	Handle(*methods.BotExt, *history.History, *types.Update)
}

// UpdateRedEnvelope 更新红包信息
func UpdateRedEnvelope(bot *methods.BotExt, redEnvelope *storage.RedEnvelope, received uint32) {
	if !redEnvelope.Active {
		return
	}
	reply := tr(0, "lng_chat_welcome")
	typ := redEnvelopesTypeToString(redEnvelope.Lucky)
	amount := fmt.Sprintf("%.2f", float64(redEnvelope.Amount)/100.0)
	if !redEnvelope.Lucky {
		amount = fmt.Sprintf("%.2f", float64(redEnvelope.Amount*redEnvelope.Number)/100.0)
	}

	reply = fmt.Sprintf(reply, typ, received, redEnvelope.Number, amount,
		storage.GetAsset(redEnvelope.Asset), redEnvelope.SenderName,
		redEnvelope.Memo, getAd(bot.ID), bot.UserName, redEnvelope.ID, bot.UserName, redEnvelope.ID)
	handler := storage.RedEnvelopeStorage{}
	if handler.IsExpired(redEnvelope.ID) {
		menus := [...]methods.InlineKeyboardButton{
			methods.InlineKeyboardButton{Text: tr(0, "lng_chat_expired"), CallbackData: "expired"},
		}
		bot.EditReplyMarkup(redEnvelope.GroupID, redEnvelope.MessageID, reply, true,
			methods.MakeInlineKeyboardMarkup(menus[:], 1))
	} else if received == redEnvelope.Number {
		menus := [...]methods.InlineKeyboardButton{
			methods.InlineKeyboardButton{Text: tr(0, "lng_chat_finished"), CallbackData: "removed"},
		}
		bot.EditReplyMarkup(redEnvelope.GroupID, redEnvelope.MessageID, reply, true,
			methods.MakeInlineKeyboardMarkup(menus[:], 1))
	} else {
		data := strconv.FormatUint(redEnvelope.ID, 10)
		menus := [...]methods.InlineKeyboardButton{
			methods.InlineKeyboardButton{Text: tr(0, "lng_chat_receive"), CallbackData: data},
		}
		bot.EditReplyMarkup(redEnvelope.GroupID, redEnvelope.MessageID, reply, true,
			methods.MakeInlineKeyboardMarkup(menus[:], 1))
	}
}

// MainMenu 主菜单
type MainMenu struct {
}

// Handle 消息处理
func (menu *MainMenu) Handle(bot *methods.BotExt, r *history.History, update *types.Update) {
	// 处理发送红包
	if update.Message != nil {
		menu.handleSendRedEnvelopes(bot, update.Message)
		return
	}

	// 处理领取红包
	if update.CallbackQuery != nil {
		menu.handleReceiveEnvelopes(bot, update.CallbackQuery)
		return
	}
}

// 消息路由
func (menu *MainMenu) route(bot *methods.BotExt, query *types.CallbackQuery) Handler {
	return nil
}

// 红包类型转字符串
func redEnvelopesTypeToString(isLucky bool) string {
	if isLucky {
		return tr(0, "lng_priv_give_rand")
	}
	return tr(0, "lng_priv_give_equal")
}

// 处理发送红包
func (menu *MainMenu) handleSendRedEnvelopes(bot *methods.BotExt, message *types.Message) {
	// 获取参数
	fromID := message.From.ID
	result := strings.Split(message.Text, " ")
	start := fmt.Sprintf("/start@%s", bot.UserName)
	if len(result) != 2 || result[0] != start {
		return
	}

	id, err := strconv.ParseUint(result[1], 10, 64)
	if err != nil {
		return
	}

	// 获取红包信息
	handler := storage.RedEnvelopeStorage{}
	redEnvelope, received, err := handler.GetRedEnvelope(id)
	if err != nil {
		logger.Errorf("Failed to get red envelope, %v", err)
		return
	}

	// 检查重复激活
	if redEnvelope.Active {
		return
	}

	// 检查红包过期
	now := time.Now().Unix()
	dynamicCfg := config.GetDynamic()
	if now-redEnvelope.Timestamp >= dynamicCfg.RedEnvelopeExpire {
		return
	}

	// 生成菜单列表
	data := strconv.FormatUint(redEnvelope.ID, 10)
	menus := [...]methods.InlineKeyboardButton{
		methods.InlineKeyboardButton{Text: tr(0, "lng_chat_receive"), CallbackData: data},
	}

	// 回复红包信息
	reply := tr(0, "lng_chat_welcome")
	typ := redEnvelopesTypeToString(redEnvelope.Lucky)
	amount := fmt.Sprintf("%.2f", float64(redEnvelope.Amount)/100.0)
	if !redEnvelope.Lucky {
		amount = fmt.Sprintf("%.2f", float64(redEnvelope.Amount*redEnvelope.Number)/100.0)
	}
	reply = fmt.Sprintf(reply, typ, received, redEnvelope.Number, amount,
		storage.GetAsset(redEnvelope.Asset), redEnvelope.SenderName,
		redEnvelope.Memo, getAd(bot.ID), bot.UserName, redEnvelope.ID, bot.UserName, redEnvelope.ID)
	markup := methods.MakeInlineKeyboardMarkup(menus[:], 1)
	message, err = bot.SendMessage(message.Chat.ID, reply, true, markup)
	if err != nil {
		logger.Errorf("Failed to send red envelope info, %v", err)
		return
	}

	// 激活红包
	err = handler.ActiveRedEnvelope(id, fromID, message.Chat.ID, message.MessageID)
	if err != nil {
		logger.Errorf("Failed to active red envelope, %v", err)
		return
	}
}

// 处理红包错误
func (menu *MainMenu) handleReceiveError(bot *methods.BotExt, query *types.CallbackQuery,
	id uint64, err error) {

	// 没有红包
	fromID := query.From.ID
	if err == storage.ErrNoBucket {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_invalid_id"), false, "", 0)
		return
	}

	// 没有激活
	if err == storage.ErrNotActivated {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_not_activated"), false, "", 0)
		return
	}

	// 领完了
	if err == storage.ErrNothingLeft {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_nothing_left"), false, "", 0)
		return
	}

	// 重复领取
	if err == storage.ErrRepeatReceive {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_repeat_receive"), false, "", 0)
		return
	}

	// 红包过期
	if err == storage.ErrRedEnvelopedExpired {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_expired"), false, "", 0)
		return
	}

	logger.Errorf("Failed to receive red envelope, id: %d, user_id: %d, %v",
		id, fromID, err)
	bot.AnswerCallbackQuery(query, tr(0, "lng_chat_receive_error"), false, "", 0)
}

// 处理领取红包
func (menu *MainMenu) handleReceiveEnvelopes(bot *methods.BotExt, query *types.CallbackQuery) {
	// 是否过期
	fromID := query.From.ID
	if query.Data == "expired" {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_expired_say"), false, "", 0)
		return
	}

	// 是否结束
	if query.Data == "removed" {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_nothing_left"), false, "", 0)
		return
	}

	// 获取红包ID
	id, err := strconv.ParseUint(query.Data, 10, 64)
	if err != nil {
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_invalid_id"), false, "", 0)
		return
	}

	// 执行领取红包
	handler := storage.RedEnvelopeStorage{}
	value, _, err := handler.ReceiveRedEnvelope(id, 0, query.From.FirstName)
	if err != nil {
		menu.handleReceiveError(bot, query, id, err)
		return
	}
	logger.Warnf("Receive red envelope, id: %d, user_id: %d, value: %d", id, fromID, value)

	// 获取红包信息
	redEnvelope, received, err := handler.GetRedEnvelope(id)
	if err != nil {
		logger.Errorf("Failed to get red envelope, %v", err)
		bot.AnswerCallbackQuery(query, tr(0, "lng_chat_receive_error"), false, "", 0)
		return
	}

	// 更新资产信息
	assetHandler := storage.AssetStorage{}
	err = assetHandler.TransferFrozenAsset(redEnvelope.SenderID, fromID,
		redEnvelope.Asset, uint32(value))
	if err != nil {
		logger.Fatalf("Failed to transfer frozen asset, from: %d, to: %d, asset: %s, amount: %d, %v",
			redEnvelope.SenderID, fromID, redEnvelope.Asset, value, err)
		return
	}

	// 更新聊天红包
	UpdateRedEnvelope(bot, redEnvelope, received)

	// 记录操作历史
	desc := fmt.Sprintf("您领取了%s(*%d*)发放的红包(id: *%d*), 获得*%.2f* *%s*", redEnvelope.SenderName, redEnvelope.SenderID,
		redEnvelope.ID, float64(redEnvelope.Amount)/100.0, redEnvelope.Asset)
	models.InsertHistory(fromID, desc)

	// 回复领取信息
	reply := tr(0, "lng_chat_receive_success")
	amount := fmt.Sprintf("%.2f", float64(value)/100.0)
	reply = fmt.Sprintf(reply, query.From.FirstName, fromID, amount,
		storage.GetAsset(redEnvelope.Asset))
	bot.ReplyMessage(query.Message, reply, true, nil)
	bot.AnswerCallbackQuery(query, tr(0, "lng_chat_receive_success_answer"), false, "", 0)

	// 回复领完消息
	if received == redEnvelope.Number {
		reply = tr(0, "lng_chat_receive_gameover")
		minRedEnvelope, maxRedEnvelope, err := handler.GetTwoTxtremes(id)
		if err == nil && redEnvelope.Number > 1 && redEnvelope.Lucky {
			body := tr(0, "lng_chat_receive_two_txtremes")
			minValue := fmt.Sprintf("%.2f", float64(minRedEnvelope.Value)/100.0)
			maxValue := fmt.Sprintf("%.2f", float64(maxRedEnvelope.Value)/100.0)
			body = fmt.Sprintf(body, maxRedEnvelope.User.FirstName, maxRedEnvelope.User.UserID, maxValue,
				storage.GetAsset(redEnvelope.Asset), minRedEnvelope.User.FirstName, minRedEnvelope.User.UserID,
				minValue, storage.GetAsset(redEnvelope.Asset))
			reply = reply + "\n\n" + body
		}
		bot.ReplyMessage(query.Message, reply, true, nil)
	}
}
