// got manhattan formula/inspiration for smart moves from https://www.geeksforgeeks.org/maximum-manhattan-distance-between-a-distinct-pair-from-n-coordinates/
// got best choice algorithm inspiration from https://www.geeksforgeeks.org/policemen-catch-thieves/
package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Board struct {
	n, m   int
	Police [2]int
	Thief  [2]int
	moves  int
}

func gameControl(board *Board, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	prevpolice := board.Police
	prevthief := board.Thief
	for {
		select {
		case msg := <-ch:
			switch msg {
			case 1:
				if prevpolice != board.Police || prevthief != board.Thief {
					fmt.Println("The Police Position = ( ", board.Police[0], board.Police[1], ") ", "The thief psotion = ( ", board.Thief[0], board.Thief[1], ") ")
					prevpolice = board.Police
					prevthief = board.Thief
				}
			case -1:

				fmt.Println("The Police caught the Thief at ( ", board.Police[0], board.Police[1], ") ", "and won the game")
				return
			case 2:
				fmt.Println("The Thief escaoed and won the game.")
				return

			}
		case <-time.After(time.Second * 2):
			if board.moves <= 0 {
				fmt.Println("The police ran out of moves and the thief won the game.")
				return
			}
		}
	}
}

func manhattanDist(x1, y1, x2, y2 int) int {
	return int(math.Abs(float64(x1-x2))) + int(math.Abs(float64(y1-y2)))
}

func movePolice(board *Board, ch chan int) {

	for board.moves > 0 {
		distance := manhattanDist(board.Police[0], board.Police[1], board.Thief[0], board.Thief[1])

		// algorithm to make best decision
		if board.Police[0] < board.Thief[0] {
			board.Police[0]++
		} else if board.Police[0] > board.Thief[0] {
			board.Police[0]--
		}

		if board.Police[1] < board.Thief[1] {
			board.Police[1]++
		} else if board.Police[1] > board.Thief[1] {
			board.Police[1]--
		}
		newDistance := manhattanDist(board.Police[0], board.Police[1], board.Thief[0], board.Thief[1])
		if newDistance > distance {
			board.Police[0] -= board.Police[0] - board.Thief[0]
			board.Police[1] -= board.Police[1] - board.Thief[1]
		}

		ch <- 1
		time.Sleep(time.Millisecond * 500)
	}
}
func moveThief(board *Board, ch chan int) {
	for {
		//algorithm to make best decisions
		if board.Police[0] < board.Thief[0] {
			board.Thief[0]--
		} else if board.Police[0] > board.Thief[0] {
			board.Thief[0]++
		}

		if board.Police[1] < board.Thief[1] {
			board.Thief[1]--
		} else if board.Police[1] > board.Thief[1] {
			board.Thief[1]++
		}

		if board.Thief == board.Police {
			ch <- -1
			return
		}

		if board.Thief == [2]int{0, 0} {
			ch <- 2
			return
		}
		ch <- 1
		time.Sleep(time.Millisecond * 500)
	}
}
func validPostion(x, y, n, m int) bool {
	return x >= 0 && x < n && y >= 0 && y < m
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {

	//random board dimensions
	n := rand.Intn(191) + 10
	m := rand.Intn(191) + 10
	//random moves police
	maxMoves := 10 * max(n, m)
	s := rand.Intn(maxMoves) + 2*max(n, m)

	//create board
	board := &Board{
		n:      n,
		m:      m,
		Police: [2]int{0, 0},
		Thief:  [2]int{n - 1, m - 1},
		moves:  s,
	}
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go movePolice(board, ch)

	wg.Add(1)
	go moveThief(board, ch)

	wg.Add(1)
	go gameControl(board, ch, &wg)
	wg.Wait()
}
