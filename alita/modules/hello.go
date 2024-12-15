package modules

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/davecgh/go-spew/spew"
	"github.com/divideprojects/Alita_Robot/alita/utils/helpers"
	log "github.com/sirupsen/logrus"
)

// Helper function to inspect struct fields
func inspectStruct(v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	log.Infof("Type: %s", t.Name())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		log.Infof("  Field: %s, Type: %s", field.Name, field.Type)
	}
}

var helloModule = moduleStruct{
	moduleName: "Hello",
}

func (moduleStruct) hello(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	user := ctx.EffectiveSender.User

	// Method 1: Basic structured logging
	log.WithFields(log.Fields{
		"user_id": ctx.EffectiveSender.User.Id,
		"chat_id": ctx.EffectiveChat.Id,
		"command": "hello",
	}).Info("Command received")

	// Method 2: JSON pretty print
	senderJSON, _ := json.MarshalIndent(ctx.EffectiveSender, "", "  ")
	log.Info("Sender structure:\n", string(senderJSON))

	// Method 3: Detailed structure dump
	log.Info("Full context dump:")
	spew.Dump(ctx.EffectiveSender)

	// Method 4: Verbose format printing
	log.Infof("Message verbose: %+v", msg)

	// Method 5: Using reflection
	log.Info("Inspecting EffectiveSender structure:")
	inspectStruct(ctx.EffectiveSender)

	// Original functionality
	replyText := fmt.Sprintf("👋 Hello %s! I'm Alita, nice to meet you!", user.FirstName)

	_, err := msg.Reply(b, replyText, helpers.Shtml())
	if err != nil {
		log.Error("Error sending hello message: ", err)
		return err
	}

	return ext.EndGroups
}

func LoadHello(dispatcher *ext.Dispatcher) {
	// Register the module
	HelpModule.AbleMap.Store("Hello", true)

	// Add command handler
	dispatcher.AddHandler(handlers.NewCommand("hello", helloModule.hello))
	log.Info("Hello module loaded successfully")
}
