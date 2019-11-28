package server

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

var questionList = [...]string{"apple juice", "orange juice", "dallas steak", "london", "berlin", "istanbul"}

// Player is model that contains all player data
type Player struct {
	CurrentGame   *Game
	Games         []Game
	IsInitialized bool
}

// Game is model for singular game data
type Game struct {
	Status       string // o:ongoing, f:finished
	Question     []QuestionItem
	MistakesLeft int
	Guesses      []Guess
	StartedAt    int64
	FinishedAt   int64
	IsWon        bool
}

// Guess is object for every single user guess
type Guess struct {
	Character string
	IsCorrect bool
}

// QuestionItem is used for every single character in question
type QuestionItem struct {
	Character string
	Guessed   bool
}

func newPlayer() Player {
	return Player{
		IsInitialized: true,
	}
}

func (p *Player) newGame() {
	rand.Seed(time.Now().Unix())

	var items []QuestionItem
	for _, rune := range questionList[rand.Intn(len(questionList))] {
		items = append(items, QuestionItem{
			Character: string(rune),
			Guessed:   unicode.IsSpace(rune),
		})
	}

	game := Game{
		Status:       "o",
		Question:     items,
		MistakesLeft: 10, // One can make 10 mistakes in classic Hangman
		StartedAt:    time.Now().Unix(),
		IsWon:        false,
	}

	p.completeGame()

	p.Games = append(p.Games, game)
	p.CurrentGame = &p.Games[len(p.Games)-1]
}

func (g *Game) showQuestion() string {
	var str strings.Builder

	for _, value := range g.Question {
		if value.Guessed {
			str.WriteString(value.Character)
		} else {
			str.WriteString("_")
		}
	}

	return str.String()
}

func (g *Game) printGameStatus() string {
	return fmt.Sprintf("Question: %v | Mistakes left: %v", g.showQuestion(), g.MistakesLeft)
}

func (g *Game) isFinished() bool {
	if g.MistakesLeft == 0 {
		return true
	}

	for _, value := range g.Question {
		if !value.Guessed {
			return false
		}
	}

	return true
}

/*
0: Already guessed before
1: Correct
2: Incorrect
*/
func (g *Game) guess(char string) int {
	for _, value := range g.Guesses {
		if value.Character == char {
			return 0
		}
	}

	ret := 2
	correct := false

	for key, value := range g.Question {
		if value.Character == char && !value.Guessed {
			g.Question[key].Guessed = true

			ret = 1
			correct = true
		}
	}

	if !correct {
		g.MistakesLeft--
	}

	g.Guesses = append(g.Guesses, Guess{
		Character: char,
		IsCorrect: correct,
	})

	return ret
}

func (p *Player) completeGame() (string, error) {
	if p.CurrentGame != nil {
		isWon := true

		if p.CurrentGame.MistakesLeft > 0 {
			for _, value := range p.CurrentGame.Question {
				if !value.Guessed {
					isWon = false

					break
				}
			}
		} else {
			isWon = false
		}

		p.CurrentGame.Status = "f"
		p.CurrentGame.FinishedAt = time.Now().Unix()
		p.CurrentGame.IsWon = isWon

		status := "You lost.."

		if isWon {
			status = "You win!"
		}
		ret := fmt.Sprintf("Game result: %v\n%s\nStart Over?", p.CurrentGame.showQuestion(), status)

		p.CurrentGame = nil

		return ret, nil
	}

	return "", errors.New("There is no game to finish")
}

// ListGames returns player's game archive
func (p *Player) ListGames() string {
	var str strings.Builder

	for _, value := range p.Games {
		str.WriteString(fmt.Sprintf("Question: %v, Status: %v, StartedAt: %v, FinishedAt: %v\n", value.showQuestion(), value.IsWon, value.StartedAt, value.FinishedAt))
	}

	return str.String()
}
