package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	var part1_p1Total, part1_p2Total int = 0, 0
	var part2_p1Total, part2_p2Total int = 0, 0

	reader := bufio.NewReader(os.Stdin)
	for {
		str, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(string(str))

		g := NewRpsRound(line)

		part1_p1Score, part1_p2Score := g.Score("Part1")
		part1_p1Total += part1_p1Score
		part1_p2Total += part1_p2Score

		part2_p1Score, part2_p2Score := g.Score("Part2")
		part2_p1Total += part2_p1Score
		part2_p2Total += part2_p2Score
	}

	fmt.Printf("Part 1 - Player 1 total score: %v\n", part1_p1Total)
	fmt.Printf("Part 1 - Player 2 total score: %v\n", part1_p2Total)
	fmt.Println("")
	fmt.Printf("Part 2 - Player 1 total score: %v\n", part2_p1Total)
	fmt.Printf("Part 2 - Player 2 total score: %v\n", part2_p2Total)
}

var ValMap = map[string]int{
	"Rock": 1, "Paper": 2, "Scissors": 3,
}

var PlayMap = map[string]string{
	"A": "Rock", "B": "Paper", "C": "Scissors",
	"X": "Rock", "Y": "Paper", "Z": "Scissors", // Only valid for part 1
}

var StratMap = map[string]string{
	"X": "lose", "Y": "draw", "Z": "win",
}

type Rps struct {
	Player1 string
	Player2 string
	Game    string
}

func NewRpsRound(gameString string) *Rps {
	p := strings.Split(gameString, " ")
	return &Rps{p[0], p[1], gameString}
}

func (g *Rps) P1Play() string { return PlayMap[g.Player1] }

func (g *Rps) P2Play() string {
	switch StratMap[g.Player2] {
	case "lose":
		switch g.P1Play() {
		case "Rock":
			return "Scissors"
		case "Paper":
			return "Rock"
		case "Scissors":
			return "Paper"
		}
	case "win":
		switch g.P1Play() {
		case "Rock":
			return "Paper"
		case "Paper":
			return "Scissors"
		case "Scissors":
			return "Rock"
		}
	case "draw":
	default:
		break
	}
	return g.P1Play()
}

func (g *Rps) Score(puzzlePart string) (int, int) {
	var draw, win int = 3, 6

	p1Play := g.P1Play()
	p2Play := g.P2Play()

	if puzzlePart == "Part1" {
		p2Play = PlayMap[g.Player2]
	}

	if p1Play == p2Play {
		return ValMap[p1Play] + draw, ValMap[p2Play] + draw
	}
	if p1Play == "Rock" && p2Play == "Scissors" {
		return ValMap[p1Play] + win, ValMap[p2Play]
	}
	if p1Play == "Rock" && p2Play == "Paper" {
		return ValMap[p1Play], ValMap[p2Play] + win
	}
	if p1Play == "Paper" && p2Play == "Scissors" {
		return ValMap[p1Play], ValMap[p2Play] + win
	}
	if p1Play == "Paper" && p2Play == "Rock" {
		return ValMap[p1Play] + win, ValMap[p2Play]
	}
	if p1Play == "Scissors" && p2Play == "Paper" {
		return ValMap[p1Play] + win, ValMap[p2Play]
	}
	if p1Play == "Scissors" && p2Play == "Rock" {
		return ValMap[p1Play], ValMap[p2Play] + win
	}
	return 0, 0
}
