package modules

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/divideprojects/Alita_Robot/alita/utils/helpers"
	log "github.com/sirupsen/logrus"
)

var helloModule = moduleStruct{
	moduleName: "Hello",
}

func (moduleStruct) hello(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	user := ctx.EffectiveSender.User

	logger := log.WithFields(log.Fields{
		"module":     "Hello",
		"command":    "hello",
		"user_id":    user.Id,
		"user_name":  user.FirstName,
		"chat_id":    ctx.EffectiveChat.Id,
		"chat_type":  ctx.EffectiveChat.Type,
		"message_id": msg.MessageId,
	})

	logger.Info("Processing hello command")

	// Create inline keyboard
	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "👋 Say Hi Back",
					CallbackData: "hello_hi_back",
				},
				{
					Text: "Bot Info",
					Url:  "https://t.me/" + b.Username,
				},
			},
		},
	}

	// Send welcome message with keyboard
	replyText := fmt.Sprintf("👋 Hello %s! I'm Alita, nice to meet you!\n\nTry out the buttons below!", user.FirstName)

	sentMsg, err := msg.Reply(b, replyText, &gotgbot.SendMessageOpts{
		ParseMode:   helpers.HTML,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to send hello message")
		return err
	}

	logger.WithFields(log.Fields{
		"sent_message_id": sentMsg.MessageId,
	}).Debug("Successfully sent hello message")

	return ext.EndGroups
}

func (moduleStruct) helloCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	user := query.From

	// Initialize logger with basic fields
	logger := log.WithFields(log.Fields{
		"module":    "Hello",
		"callback":  "hello_hi_back",
		"user_id":   user.Id,
		"user_name": user.FirstName,
	})

	// Handle message information if available
	if msg, ok := query.Message.(*gotgbot.Message); ok {
		logger = logger.WithFields(log.Fields{
			"message_id": msg.MessageId,
			"chat_id":    msg.Chat.Id,
			"chat_type":  msg.Chat.Type,
		})
	}

	logger.Info("Processing hello callback")

	// Answer callback query
	_, err := query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Hello there! 👋",
	})
	if err != nil {
		logger.WithError(err).Error("Failed to answer callback query")
		return err
	}

	// Edit original message if accessible
	if msg, ok := query.Message.(*gotgbot.Message); ok {
		editedMsg, _, err := msg.EditText(b,
			fmt.Sprintf("Nice to meet you %s! How can I help you today?", user.FirstName),
			&gotgbot.EditMessageTextOpts{
				ParseMode: helpers.HTML,
			},
		)
		if err != nil {
			logger.WithError(err).Error("Failed to edit message")
			return err
		}

		logger.WithFields(log.Fields{
			"edited_message_id": editedMsg.MessageId,
		}).Debug("Successfully processed callback")
	} else {
		logger.Debug("Message is not accessible, skipping edit")
	}

	return ext.EndGroups
}

func LoadHello(dispatcher *ext.Dispatcher) {
	log.WithFields(log.Fields{
		"module": "Hello",
	}).Info("Loading hello module")

	// Register the module
	HelpModule.AbleMap.Store("Hello", true)

	// Add command handler
	dispatcher.AddHandler(handlers.NewCommand("hello", helloModule.hello))

	// Add callback handler
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("hello_hi_back"), helloModule.helloCallback))

	log.WithFields(log.Fields{
		"module":   "Hello",
		"handlers": []string{"hello", "hello_hi_back"},
	}).Info("Hello module loaded successfully")
}
