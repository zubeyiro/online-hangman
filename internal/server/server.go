package server

import (
	"bufio"
	"fmt"
	"github.com/zubeyiro/online-hangman/internal/message"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

// Server is general server struct
type Server struct {
	listener     *net.TCPListener
	connections  map[net.Conn]uint64 // player_id
	players      map[uint64]*Player
	mtx          *sync.RWMutex
	lastPlayerID uint64
	shutdown     chan bool
}

// New creates new server instance
func New() *Server {
	server := Server{
		connections: make(map[net.Conn]uint64),
		players:     make(map[uint64]*Player),
		mtx:         &sync.RWMutex{},
		shutdown:    make(chan bool),
	}

	return &server
}

// Start function starts listening for clients
func (server *Server) Start() error {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 5454})

	if err != nil {
		return err
	}
	server.listener = listener

	fmt.Println("# Server started listening")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("# Error while new connection: ", err)
		} else {
			go server.handleConnection(conn)
		}
	}
}

// Stop terminates server
func (server *Server) Stop() error {
	server.shutdown <- true

	return server.listener.Close()
}

func (server *Server) newConnection(conn net.Conn, playerID uint64) uint64 {
	server.mtx.Lock()
	defer server.mtx.Unlock()

	var pID uint64

	if playerID == 0 {
		pID = atomic.AddUint64(&server.lastPlayerID, 1)
	} else {
		pID = playerID
	}

	server.connections[conn] = pID

	return pID
}

func (server *Server) removeConnection(conn net.Conn) {
	server.mtx.Lock()
	defer server.mtx.Unlock()

	conn.Close()
	delete(server.connections, conn)
}

func (server *Server) handleConnection(conn net.Conn) {
	var playerID uint64
	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	error := make(chan error)

	go func() {
		for {
			buffer, err := stream.Reader.ReadString('\n')

			if err != nil {
				error <- err
				return
			}

			if msg, err := message.Parse(buffer); err != nil {
				error <- err
			} else {
				if playerID > 0 {
					player := server.players[playerID]
					switch msg.Command {
					case message.CliStartGame:
						server.mtx.Lock()

						player.newGame()
						conn.Write([]byte(fmt.Sprintf("%v\n", player.CurrentGame.printGameStatus())))

						server.mtx.Unlock()
					case message.CliGuess:
						server.mtx.Lock()

						if player.CurrentGame != nil {
							if len(msg.Payload) > 1 {
								conn.Write([]byte("Your guess has to be 1 character\n"))
							} else {
								guessResult := player.CurrentGame.guess(msg.Payload)

								if guessResult == 0 {
									conn.Write([]byte("You have guessed this character before\n"))
								} else {
									if guessResult == 1 {
										conn.Write([]byte("Correct!\n"))
									} else {
										conn.Write([]byte("Incorrect!\n"))
									}

									if player.CurrentGame.isFinished() {
										if result, err := player.completeGame(); err != nil {
											conn.Write([]byte(fmt.Sprintf("%v\n", err)))
										} else {
											conn.Write([]byte(fmt.Sprintf("%v\n", result)))
										}
									} else {
										conn.Write([]byte(fmt.Sprintf("%v\n", player.CurrentGame.printGameStatus())))
									}
								}
							}
						} else {
							conn.Write([]byte("You should start game first\n"))
						}

						server.mtx.Unlock()
					case message.CliShowQuestion:
						server.mtx.Lock()

						if player.CurrentGame != nil {
							conn.Write([]byte(fmt.Sprintf("%v\n", player.CurrentGame.printGameStatus())))
						} else {
							conn.Write([]byte("You should start game first\n"))
						}

						server.mtx.Unlock()
					case message.CliListGames:
						server.mtx.Lock()

						conn.Write([]byte(fmt.Sprintf("%v\n", player.ListGames())))

						server.mtx.Unlock()
					}

				} else {
					if msg.Command == message.CliLogin {
						if id, err := strconv.ParseUint(msg.Payload, 10, 64); err != nil {
							error <- err
						} else {
							server.mtx.Lock()
							if server.players[id] != nil {
								server.mtx.Unlock()
								playerID = server.newConnection(conn, id)
								conn.Write([]byte("Welcome back!\n"))
							} else {
								server.mtx.Unlock()
								playerID = server.newConnection(conn, 0)
								newPlayer := newPlayer()
								server.mtx.Lock()
								server.players[playerID] = &newPlayer
								server.mtx.Unlock()
								conn.Write([]byte("Welcome!\n"))
							}
						}
					} else {
						conn.Write([]byte("Login first\n"))
					}
				}
			}
		}
	}()

	for {
		select {
		case <-server.shutdown:
			break
		case err := <-error:
			fmt.Println("# Error on server: ", err)
			conn.Write([]byte(fmt.Sprintf("Error on server: %v", err)))
			server.removeConnection(conn)
			return
		}
	}
}
