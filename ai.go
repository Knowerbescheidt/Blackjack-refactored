package blackjack

//minute 9 ai interfrace
import (
	"fmt"

	deck "github.com/Knowerbescheidt/Deck-of-cards"
)

type AI interface {
	Results(hand [][]deck.Card, dealer []deck.Card)
	Play(hand []deck.Card, dealer deck.Card) Move
	Bet() int
}

type dealerAI struct{}

func (ai dealerAI) Bet() int {
	return 1
}

func (ai *dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore < 16 || dScore == 17 && Soft(hand...) {
		return MoveHit
	} else {
		return MoveStand
	}
}

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	//noop
}

type HumanAI struct{}

func (ai *HumanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do?(h)it or (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return Hit
		case "s":
			return Stand
		default:
			fmt.Println("Invalid Option", input)
		}
	}
}

func (ai *HumanAI) Bet() int {
	return 1
}
