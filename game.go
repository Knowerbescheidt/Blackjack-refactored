package blackjack

import (
	"fmt"

	deck "github.com/Knowerbescheidt/Deck-of-cards"
)

//minute 31 refactor
type state int8

func New() Game {
	return Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
}

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Game struct {
	//unexported fields
	deck     []deck.Card
	state    state
	player   []deck.Card
	dealer   []deck.Card
	dealerAI AI
	balance  int
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("A valid Gamestate can not be found")
	}
}

type Move func(*Game)

func MoveHit(g *Game) {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func MoveStand(g *Game) {
	g.state++
}

func (g *Game) Play(ai AI) int {
	g.deck = deck.New(deck.Deck(3), deck.Shuffle)

	//presumably 10 games
	for i := 0; i < 10; i++ {
		deal(g)

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(g.player, g.dealer[0])
			move(g)
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}
		endHand(g, ai)
	}
	return 0
}

func deal(g *Game) {
	g.player = make([]deck.Card, 0, 7)
	g.dealer = make([]deck.Card, 0, 7)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = g.deck[0], g.deck[1:]
		g.player = append(g.player, card)
		card, g.deck = g.deck[0], g.deck[1:]
		g.dealer = append(g.dealer, card)
	}
	g.state = statePlayerTurn
}

func Score(h ...deck.Card) int {
	//does not copy the h
	minScore := minScore(h...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func Soft(h ...deck.Card) bool {
	min := minScore(h...)
	sc := Score(h...)
	return min != sc
}

func minScore(h ...deck.Card) int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	switch {
	case pScore > 21:
		fmt.Println("You busted, and loose, Looser")
		g.balance--
	case dScore > 21:
		fmt.Println("Dealer busted, and loose, Looser")
		g.balance++
	case pScore > dScore:
		fmt.Println("You win Congrats")
		g.balance++
	case dScore > pScore:
		fmt.Println("You loose...Try Again!")
		g.balance--
	case dScore == pScore:
		fmt.Println("Draw")
	}
	fmt.Println()
	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.player = nil
}
