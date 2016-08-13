package main

import (
	"fmt"
	"sync"
	"github.com/IonPostglacial/othello"
)

type result struct {
	score int
	key othello.Position
	move []othello.Position
}

func main() {
	board := new(othello.Board)
	board.SetCell(3, 3, othello.WHITE)
	board.SetCell(4, 4, othello.WHITE)
	board.SetCell(3, 4, othello.BLACK)
	board.SetCell(4, 3, othello.BLACK)
	fmt.Println(board)

	currentPlayer := othello.BLACK

	for {
		possibleMoves := othello.PossibleMoves(board, currentPlayer)
		if len(possibleMoves) == 0 {
			fmt.Println(currentPlayer, ": ", board.Count(currentPlayer))
			fmt.Println(currentPlayer.Opponent(), ": ", board.Count(currentPlayer.Opponent()))
			break
		}
		var (
			bestMove  []othello.Position
			bestScore int
			bestKey othello.Position
			wg sync.WaitGroup
		)
		done := make(chan bool)
		results := make(chan result)
		bestMoveChan := make(chan []othello.Position)
		wg.Add(len(possibleMoves))
		go func() {
			wg.Wait()
			done <- true
		}()
		for key, move := range possibleMoves {
			go func(move []othello.Position) {
				results <- result {
					score: othello.EvaluateMove(move, board, currentPlayer, currentPlayer),
					key: key,
					move: move,
				}
				wg.Done()
			}(move)
		}
		go func() {
			for {
				select {
				case <-done:
					bestMoveChan <- bestMove
					return
				case result := <-results:
					if result.score >= bestScore {
						bestKey = result.key
						bestScore = result.score
						bestMove = result.move
					}
				}
			}
		}()
		for _, position := range <-bestMoveChan {
			board.SetCell(position.H, position.V, currentPlayer)
		}
		fmt.Println(currentPlayer, " played ", bestKey, "[", bestScore, "]: ", bestMove)
		fmt.Println(board)
		currentPlayer = currentPlayer.Opponent()
	}
}
