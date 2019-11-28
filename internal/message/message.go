package message

import (
	"errors"
	"strings"
)

var CliLogin = "cl"
var CliStartGame = "cs"
var CliGuess = "cg"
var CliShowQuestion = "csq"
var CliListGames = "csl"

var commands = [...]string{CliLogin, CliStartGame, CliGuess, CliShowQuestion, CliListGames}

// Message is message object thats used between server and client
type Message struct {
	Command string
	Payload string
}

// Parse imports raw socket message and returns as message format is correctly sent
func Parse(raw string) (Message, error) {
	raw = strings.TrimSpace(raw)

	if strings.Contains(raw, "-") {
		rawArr := strings.Split(raw, "-")
		cmd := rawArr[0]
		payload := rawArr[1]

		if commandExists(cmd) {
			return Message{
				Command: cmd,
				Payload: payload,
			}, nil
		}

		return Message{}, errors.New("Incorrect command")
	}

	return Message{}, errors.New("Incorrect command")
}

func commandExists(cmd string) bool {
	for _, value := range commands {
		if value == cmd {
			return true
		}
	}

	return false
}
