package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

// Send any text message to the bot after the bot has been started

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "kaguya", bot.MatchTypePrefix, aiHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "weweng", bot.MatchTypePrefix, wewengBading)
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Say kaguya and I'll answer.\nExample: kaguya who invented calculus?",
	})
}

func aiHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      chatCompletion(update.Message.Text),
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("message: %v\n", message)
	}
}

func wewengBading(ctx context.Context, b *bot.Bot, update *models.Update) {
	badingResp := "As if I would be interested in programming or anything like that, but fine! Here’s a simple example of a \"Hello, World!\" program in Rust:\n\n```rust\nfn main() {\n    println!(\"Hello, world!\");\n}\n```\n\nYou just need to create a new Rust file (like `main.rs`), put that code in, and run it using `cargo run` or `rustc main.rs` followed by `./main` if you’re compiling directly. It’s not like I’m trying to help you or anything! Just don’t get too carried away, okay?"

	message, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      badingResp,
		ParseMode: models.ParseModeMarkdownV1,
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("message: %v\n", message)
	}
}
