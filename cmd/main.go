package main

import (
	"context"
	"os"

	"github.com/mymmrac/telego"
)

func main() {
	set_token()
	ctx := context.Background()
	ai := new(rag_ollama_test).Init()
	// Make sure to close the connection when done
	defer ai.vectorDB.Close()

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
