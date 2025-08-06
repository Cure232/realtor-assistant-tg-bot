package main

import (
	"context"
	"os"

	"github.com/mymmrac/telego"
)

func main() {
	set_token()
	ctx := context.Background()
	ai := new(ragLangchainTest).Init()
	//ai.vectorStore.AddDocuments(ctx, []schema.Document{{PageContent: "Высадка на Луну произошла в 2023 году. Согласно заявлению Роскосмоса."}})
	// Make sure to close the connection when done
	defer ai.vectorStore.Close()

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
	HandleUpdates(bot, ai, updates)
}

func set_token() {
	file, err := os.ReadFile("../token.txt")
	if err != nil {
		panic(err)
	}
	os.Setenv("TOKEN", string(file))
}
