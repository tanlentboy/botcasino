package privatechat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zhangpanyi/botcasino/models"

	"github.com/zhangpanyi/basebot/history"
	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// 匹配历史页数
var reMathHistoryPage *regexp.Regexp

func init() {
	var err error
	reMathHistoryPage, err = regexp.Compile("^/history/(|(\\d+)/)$")
	if err != nil {
		panic(err)
	}
}

// History 历史记录
type History struct {
}

// Handle 消息处理
func (history *History) Handle(bot *methods.BotExt, r *history.History, update *types.Update) {
	data := update.CallbackQuery.Data
	result := reMathHistoryPage.FindStringSubmatch(data)
	if len(result) == 3 {
		page, err := strconv.Atoi(result[2])
		if err != nil {
			history.replyHistory(bot, 0, update.CallbackQuery)
		} else {
			history.replyHistory(bot, page, update.CallbackQuery)
		}
	}
}

// 消息路由
func (history *History) route(bot *methods.BotExt, query *types.CallbackQuery) Handler {
	return nil
}

// 生成菜单列表
func (history *History) makeMenuList(fromID int64, page int) *methods.InlineKeyboardMarkup {
	priv := fmt.Sprintf("/history/%d/", page-1)
	next := fmt.Sprintf("/history/%d/", page+1)
	menus := [...]methods.InlineKeyboardButton{
		methods.InlineKeyboardButton{Text: tr(fromID, "lng_previous_page"), CallbackData: priv},
		methods.InlineKeyboardButton{Text: tr(fromID, "lng_next_page"), CallbackData: next},
		methods.InlineKeyboardButton{Text: tr(fromID, "lng_back_superior"), CallbackData: "/main/"},
	}
	return methods.MakeInlineKeyboardMarkupAuto(menus[:], 2)
}

// 生成回复内容
func (history *History) makeReplyContent(fromID int64, array []models.History, page, pagesum uint) string {
	header := fmt.Sprintf("%s (*%d*/%d)\n\n", tr(fromID, "lng_priv_history"), page, pagesum)
	if len(array) > 0 {
		lines := make([]string, 0, len(array))
		format := tr(fromID, "lng_priv_history_fmt")
		for _, his := range array {
			date := his.InsertedAt.Format("2006-01-02 03:04:05")
			lines = append(lines, fmt.Sprintf(format, date, his.Describe))
		}
		return header + strings.Join(lines, "\n\n")
	}
	return header + tr(fromID, "lng_priv_history_no_op")
}

// 回复历史记录
func (history *History) replyHistory(bot *methods.BotExt, page int, query *types.CallbackQuery) {
	// 检查页数
	if page < 1 {
		page = 1
	}

	// 查询历史
	fromID := query.From.ID
	array, pagesum, err := models.GetUserHistory(fromID, page, 5)
	if err != nil {
		logger.Warnf("Failed to query user history, %v", err)
	}
	if page > int(pagesum) {
		page = int(pagesum)
	}

	// 回复内容
	if len(array) > 0 {
		bot.AnswerCallbackQuery(query, "", false, "", 0)
	} else {
		reply := tr(fromID, "lng_priv_history_no_op")
		bot.AnswerCallbackQuery(query, reply, false, "", 0)
	}
	reply := history.makeReplyContent(fromID, array, uint(page), pagesum)
	bot.EditMessageReplyMarkup(query.Message, reply, true, history.makeMenuList(fromID, page))
}
