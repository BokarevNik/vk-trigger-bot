package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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

	cmd, err := readCommand()
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		if obj.Message.Text == "/shutdown" {
			log.Printf("/shutdown received")

			kill := exec.Command(cmd[0], cmd[1:]...)
			if out, err := kill.CombinedOutput(); err != nil {
				log.Fatalf("err: %s\noutput: %s", err.Error(), string(out))
			}

			log.Println("process killed")
		}
	})

	// Run Bots Long Poll
	log.Println("Start polling")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}

func readCommand() ([]string, error) {
	var inputCmd string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter command: ")

	if scanner.Scan() {
		inputCmd = scanner.Text()
	} else {
		return nil, fmt.Errorf("Input reading error")
	}

	cmd := strings.Fields(inputCmd)

	if len(cmd) <= 1 {
		return nil, fmt.Errorf("No arguments for command")
	}

	return cmd, nil
}
