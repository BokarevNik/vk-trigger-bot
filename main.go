package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func main() {
	vk := api.NewVK(token)

	// Initializing Long Poll
	lp, err := longpoll.NewLongPollCommunity(vk)
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		if obj.Message.Text == "/shutdown" {
			cmd := exec.Command("pkill", "zoom")
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			// log.Println("process killed")

			// b := params.NewMessagesSendBuilder()
			// b.Message("process killed")
			// b.RandomID(0)
			// b.PeerID(obj.Message.PeerID)

			// _, err := vk.MessagesSend(b.Params)
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}
	})

	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
