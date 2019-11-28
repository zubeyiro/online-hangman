package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zubeyiro/online-hangman/internal/client"
	"github.com/zubeyiro/online-hangman/internal/server"
)

func main() {
	server := server.New()
	go server.Start()
	time.Sleep(time.Second * 1) // wait for server to launch

	help := `Welcome to hangman! Here is how you can play;
login: Login as guest
login [ID]: Login with your existing account, (i.e. login-1)
start: starts the game, new game will be started if you are already playing
show: question will be prompted again
g [CHAR]: this is your move for guessing characters (i.e. g a, g b, g c)
archive: this will print your game archive
help: read this message again
	`

	cli := client.New()
	err := cli.Connect()

	if err != nil {
		fmt.Println("# Error while connecting to server: ", err)
	} else {
		clientCh := make(chan string)
		go cli.HandleIncomingMessages(clientCh)
		go func() {
			for d := range clientCh {
				fmt.Println(d)
			}
		}()

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println(help)

		for scanner.Scan() {
			input := scanner.Text()
			switch input {
			case "login":
				cli.LoginGuest()
			case "start":
				cli.StartGame()
			case "show":
				cli.ShowQuestion()
			case "archive":
				cli.ShowArchive()
			case "help":
				fmt.Println(help)
			default:
				if strings.Contains(input, " ") {
					switch strings.Split(input, " ")[0] {
					case "login":
						cli.LoginWithID(strings.Split(input, " ")[1])
					case "g":
						cli.Guess(strings.Split(input, " ")[1])
					default:
						fmt.Println("Incorrect command, use help for command list")
					}
				} else {
					fmt.Println("Incorrect command, use help for command list")
				}
			}
		}

		if scanner.Err() != nil {
			fmt.Println("# Error while reading client input")
		}

		cli.Close()
		server.Stop()
	}
}
