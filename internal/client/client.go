package client

import (
	"bufio"
	"fmt"
	"net"
)

/*
Login_with_user_id
Login_without_user_id
Start game
show question
Guess
List archive games
*/

// Client struct
type Client struct {
	conn net.Conn
	stop chan bool
}

// New creates new client
func New() *Client {
	client := Client{stop: make(chan bool)}

	return &client
}

// Connect to server
func (cli *Client) Connect() error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: 5454})

	if err != nil {
		return fmt.Errorf("# Cannot connect to server: %s", err)
	}

	cli.conn = conn

	return nil
}

// Close connection
func (cli *Client) Close() error {
	close(cli.stop)
	cli.conn.Close()

	return nil
}

// SendMessage to server
func (cli *Client) sendMessage(msg string) error {
	_, err := cli.conn.Write([]byte(fmt.Sprintf("%v\n", msg)))

	if err != nil {
		fmt.Println(fmt.Sprintf("# Error while sending message to server: %v", err))

		return err
	}

	return nil
}

// HandleIncomingMessages from server
func (cli *Client) HandleIncomingMessages(writeCh chan<- string) {
	error := make(chan error)
	r := bufio.NewReader(cli.conn)

	for {
		select {
		case <-cli.stop:
			close(error)

			return
		case err := <-error:
			writeCh <- err.Error()
			cli.Close()

			return
		default:
			msgBuffer, err := r.ReadString('\n')

			if err != nil {
				error <- err

				return
			}

			writeCh <- msgBuffer
		}
	}
}

// LoginGuest login request as guest
func (cli *Client) LoginGuest() error {
	return cli.sendMessage("cl-0")
}

// LoginWithID for existing account
func (cli *Client) LoginWithID(id string) error {
	return cli.sendMessage(fmt.Sprintf("cl-%v", id))
}

// StartGame request to server
func (cli *Client) StartGame() error {
	return cli.sendMessage("cs-0")
}

// ShowQuestion for ongoing game
func (cli *Client) ShowQuestion() error {
	return cli.sendMessage("csq-0")
}

// ShowArchive games
func (cli *Client) ShowArchive() error {
	return cli.sendMessage("csl-0")
}

// Guess character for ongoing game
func (cli *Client) Guess(char string) error {
	return cli.sendMessage(fmt.Sprintf("cg-%v", char))
}
