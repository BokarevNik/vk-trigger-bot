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

	input, err := readCommand()
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		if obj.Message.Text == "/shutdown" {
			name, err := getName(obj.Message.FromID, vk)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("/shutdown triggered by %v\n", name)

			cmd := exec.Command(input[0], input[1:]...)
			if out, err := cmd.CombinedOutput(); err != nil {
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

func getName(id int, vk *api.VK) ([]string, error) {
	sender, err := vk.UsersGet(api.Params{
		"user_ids": id,
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting user by id\nerr: %v", err)
	}

	return []string{sender[0].FirstName, sender[0].LastName}, nil
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

	if len(cmd) < 2 {
		return nil, fmt.Errorf("No arguments for command")
	}

	return cmd, nil
}
