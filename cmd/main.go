package main

import (
	"context"
	"os"

	"github.com/mymmrac/telego"
)

func main() {

	ctx := context.Background()
	ai := rag_ollama_test{}
	aib, _ := ai.Init()

	os.Setenv("TOKEN", "TOKENVALUE")
	botToken := os.Getenv("TOKEN")
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}

	cmds := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "/start", Description: "Получить ознакомительную информацию"},
			{Command: "/menu", Description: "Меню выбора инструментов"},
		},
	}
	bot.SetMyCommands(ctx, &cmds)

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)
	HandleUpdates(bot, aib, updates)

}
