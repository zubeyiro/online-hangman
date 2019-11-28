package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zubeyiro/online-hangman/internal/client"
	"github.com/zubeyiro/online-hangman/internal/server"
)

var clientCount = 1000

func TestBenchmark(t *testing.T) {
	srv := server.New()

	go func() {
		err := srv.Start()
		require.NoError(t, err)

		defer assert.NoError(t, srv.Stop())
	}()

	time.Sleep(time.Second * 1) // wait for server to launch

	var clients []*client.Client
	var clientChs []chan string
	for i := 0; i < clientCount; i++ {
		cli := client.New()
		defer func() {
			assert.NoError(t, cli.Close())
		}()

		require.NoError(t, cli.Connect())

		clientCh := make(chan string)
		clients = append(clients, cli)
		clientChs = append(clientChs, clientCh)
	}

	t.Run("error tests for commands", func(t *testing.T) {
		result := testing.Benchmark(func(b *testing.B) {
			for i := 0; i < clientCount; i++ {
				assert.NoError(b, clients[i].LoginGuest())
				assert.NoError(b, clients[i].StartGame())
				assert.NoError(b, clients[i].Guess("a"))
				assert.NoError(b, clients[i].ShowQuestion())
				assert.NoError(b, clients[i].ShowArchive())
				assert.NoError(b, clients[i].StartGame())
			}
		})
		t.Logf("error tests for commands\n%s\n", result.String())
	})
}
