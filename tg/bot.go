package tg

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wilsonce/connectly-test/initialize"
	"github.com/wilsonce/connectly-test/model"
	"github.com/wilsonce/connectly-test/service"
	"log"
	"sync"
)

type BotBuilder struct {
	buildChannel chan string
	botMap       *sync.Map
}

var botBuilderOnce sync.Once
var botBuilder *BotBuilder

func NewBotBuilder() *BotBuilder {
	botBuilderOnce.Do(func() {
		tmp := &BotBuilder{
			buildChannel: make(chan string, 10),
			botMap:       &sync.Map{},
		}
		botBuilder = tmp
	})
	return botBuilder
}

func (b *BotBuilder) Build(token string) *Bot {
	bot := NewBot()
	b.botMap.Store(token, bot)
	go bot.New(token)
	return bot
}

func (b *BotBuilder) Get(token string) (*Bot, bool) {
	tmp, ok := b.botMap.Load(token)
	if ok {
		return tmp.(*Bot), ok
	}
	return nil, false
}

func (b *BotBuilder) Remove(token string) {
	tmp, ok := b.Get(token)
	if ok {
		tmp.Cancel()
		b.botMap.Delete(token)
		tmp = nil
	}
}

func (b *BotBuilder) RemoveAll() {
	b.botMap.Range(func(key, value interface{}) bool {
		b.Remove(key.(string))
		return true
	})
}

func (b *BotBuilder) AutoBuild() {
	for {
		select {
		case token := <-b.buildChannel:
			b.Build(token)
		}
	}
}

type Bot struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBot() *Bot {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bot{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (b *Bot) NewWithWebhook(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		initialize.Logger.Error(err.Error())
		return err
	}

	bot.Debug = true
	initialize.Logger.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	//wh, _ := tgbotapi.NewWebhookWithCert("https://www.example.com:8443/"+bot.Token, "cert.pem")
	wh, _ := tgbotapi.NewWebhook("http://127.0.0.1:9999/" + bot.Token)

	_, err = bot.Request(wh)
	if err != nil {
		initialize.Logger.Error(err.Error())
		return err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		initialize.Logger.Error(err.Error())
		return err
	}

	if info.LastErrorDate != 0 {
		initialize.Logger.Info(fmt.Sprintf("Telegram callback failed: %s", info.LastErrorMessage))
		return errors.New(info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	for update := range updates {
		log.Printf("%+v\n", update)
	}

	return nil
}

func (b *Bot) New(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		initialize.Logger.Error(err.Error())
		return err
	}
	defer bot.StopReceivingUpdates()
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		select {
		case <-b.ctx.Done():
			initialize.Logger.Info(fmt.Sprintf("Bot %s exit", bot.Self.UserName))
			return nil
		default:
			if update.Message != nil { // If we got a message
				messageService := service.NewMessageService()
				messageService.Add(context.Background(), &model.WMessage{
					FromUsername:  update.Message.From.UserName,
					FromFirstName: update.Message.From.FirstName,
					FromLastName:  update.Message.From.LastName,
					Messaage:      update.Message.Text,
					ChatID:        update.Message.Chat.ID,
					BotName:       bot.Self.UserName,
				})
				initialize.Logger.Info(fmt.Sprintf("Bot give messaage username:%s, text: %s", update.Message.From.UserName, update.Message.Text))

				if update.Message.IsCommand() {
					command := update.Message.Command()
					switch command {
					case "recent_orders":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Recent orders:\r\n\r\n1. [10000000000000000000](https://www.google.com)\r\n2. [200000000000000000000](https://www.googoe.com)")
						msg.ReplyToMessageID = update.Message.MessageID
						msg.ParseMode = "markdown"
						bot.Send(msg)
					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}
		}
	}

	return nil
}

func (b *Bot) Cancel() {
	b.cancel()
}
