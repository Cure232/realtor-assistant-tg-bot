package main

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleUpdates(bot *telego.Bot, ai *ragLangchainTest, updates <-chan telego.Update, options ...th.BotHandlerOption) {
	bh, _ := th.NewBotHandler(bot, updates)
	defer func() { _ = bh.Stop() }()

	//HandleMenuButtons
	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		data := query.Data
		response := ""
		switch data {
		case "aiconsult":
			response = "Задайте свой вопрос ИИ ассистенту:"
		case "analyse":
			response = "Отправьте ссылку на объект или информацию о нём:"
		case "describe":
			response = "Отправьте ссылку на объект или информацию о нём:"
		case "autopost":
			response = "Отправьте ссылку на свой телеграм канал:"
		case "news":
			response = "Отправьте список из ссылок на Телеграм каналы из которых мы будем формировать дайджест:"
		}
		bot.SendMessage(ctx, tu.Message(
			query.Message.GetChat().ChatID(),
			response,
		))
		// Подтверждаем callback
		bot.AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return nil
	}, th.AnyCallbackQuery())

	//HandleCommands
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		text := message.Text
		switch text {
		case "/start":
			ctx.Bot().SendMessage(ctx, tu.Message(message.Chat.ChatID(), "Привет, я бот Недвижик я могу помочь тебе с... lorem ipsum"))
		case "/menu":
			greeting := "Добрый день, чем могу быть полезен?"
			keyboard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Проконсультироваться с ИИ ассистентом ❌").WithCallbackData("aiconsult"),
				),
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Получить аналитику по объекту недвижимости ❌").WithCallbackData("analyse"),
				),
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Написать описание для объекта недвижимости ❌").WithCallbackData("describe"),
				),
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Настроить автопубликацию контента для ТГ Канала ❌").WithCallbackData("autopost"),
				),
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Настроить персональный новостной дайджест ❌").WithCallbackData("news"),
				),
			)

			// Отправка сообщения с inline-клавиатурой
			bot.SendMessage(ctx, tu.Message(
				message.Chat.ChatID(),
				greeting,
			).WithReplyMarkup(keyboard))

		default:
			answer := ai.Test(ctx, text)
			bot.SendMessage(ctx, tu.Message(
				message.Chat.ChatID(),
				answer,
			))
		}
		return nil
	}, th.AnyMessageWithText())

	bh.Start()
}
