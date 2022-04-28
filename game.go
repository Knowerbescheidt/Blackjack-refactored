package blackjack

import (
	"fmt"

	deck "github.com/Knowerbescheidt/Deck-of-cards"
)

//Starting Betting
type state int8

func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.BlackjackPayout == 0 {
		opts.BlackjackPayout = 1.5
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	g.nDecks = opts.Decks
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout

	return g
}

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

type Game struct {
	//unexported fields
	nDecks          int
	nHands          int
	blackjackPayout float64

	state state
	deck  []deck.Card

	player    []deck.Card
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI AI
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

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	g.playerBet = bet
}

func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3

	//presumably 10 games
	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
			shuffled = true
		}
		bet(g, ai, shuffled)
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
	winnings := g.playerBet
	switch {
	case pScore > 21:
		fmt.Println("You busted, and loose, Looser")
		winnings *= -1
	case dScore > 21:
		fmt.Println("Dealer busted, and loose, Looser")
	case pScore > dScore:
		fmt.Println("You win Congrats")
	case dScore > pScore:
		fmt.Println("You loose...Try Again!")
		winnings *= -1
	case dScore == pScore:
		fmt.Println("Draw")
		winnings = 0
	}
	g.balance += winnings
	fmt.Println()
	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.player = nil
}
