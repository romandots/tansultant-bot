package app

import (
	"math/rand"
	"tansulbot/pkg/telegram"
	"testing"
	"time"
)

func TestSaveAndGetConversation(t *testing.T) {
	a, _ := InitApp("123")

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000)

	_, exists, _ := a.GetConversation(id)

	if exists == true {
		t.Error("Record exists")
	}

	conv := &telegram.Conversation{ChatId: id, UserId: id, State: telegram.Unauthorized}

	err := a.SaveConversation(conv)
	if err != nil {
		t.Error(err)
	}

	loadedConv, exists, err := a.GetConversation(id)
	if exists != true {
		t.Error("Record does not exist")
	}

	if loadedConv.ChatId != id {
		t.Error("Data corrupt")
	}

	loadedConv.State = telegram.Idle
	a.SaveConversation(loadedConv)

	reloadedConv, exists, err := a.GetConversation(id)
	if exists != true {
		t.Error("Record does not exist")
	}

	if reloadedConv.State != telegram.Idle {
		t.Error("Data corrupt")
	}
}
