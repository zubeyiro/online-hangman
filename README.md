# online-hangman
Online hangman game written in Go.
Clone this repository then run

`go run src/main.go`

You can play game through terminal.

## Capabilities
- User can start a new game
- User can guess a character for an ongoing game
- User is being notified about game result
- User can resume game later
- User can list games played before
- User is being notified error happens on server

## Testing
Run;

`go test -v test/benchmark_test.go`
`go test -race test/benchmark_test.go`

It's basically checks if server/client running correctly, test should be run after every update.