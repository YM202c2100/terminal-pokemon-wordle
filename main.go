package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"
)

const (
	MaxNum      int    = 292
	N           int    = 5
	BackBlack   string = "40"
	BackGreen   string = "42"
	BackMagenta string = "45"
	WordWhite   string = "37"
)

//go:embed Pokemon_jp_DP.txt
var files embed.FS

type Game struct {
	Num            int
	TurnsRemaining int
	Complete       bool
	Won            bool
	Answer         []rune
}

type Guess struct {
	Pokemon    []rune
	Background []string
}

//ランダムな整数を得た後、回答を得る
func set_answer(game *Game) {
	rand.Seed(time.Now().UnixNano())
	game.Num = rand.Intn(MaxNum) + 1
	f, err := files.Open("Pokemon_jp_DP.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)

	for count := 1; scanner.Scan(); count++ {
		if count == game.Num {
			game.Answer = []rune(scanner.Text())
			break
		}
	}
	defer f.Close()
}

//正誤判定
func (game *Game) comparing(guess *Guess) {
OuterLoop:
	for i := 0; i < N; i++ {
		if guess.Pokemon[i] == game.Answer[i] {
			guess.Background[i] = BackGreen
		} else {
			for j := 0; j < N; j++ {
				if guess.Pokemon[i] == game.Answer[j] {
					guess.Background[i] = BackMagenta
					continue OuterLoop
				}
			}
			guess.Background[i] = BackBlack
		}
	}
}

//空のボードを表示する
func (game *Game) print_empty() {
	for i := 0; i < game.TurnsRemaining; i++ {
		for j := 0; j < N; j++ {
			fmt.Print("_ ")
		}
		fmt.Println()
	}
	fmt.Println()
}

//ボードを表示させる
func (game *Game) display(guess *Guess) {
	fmt.Printf("\033[%dF\033[J", game.TurnsRemaining+3)
	for i := 0; i < N; i++ {
		fmt.Printf("\033[%sm%s\033[0m", guess.Background[i], string(guess.Pokemon[i]))
	}
	fmt.Println()
	game.print_empty()
}

func main() {
	valid := regexp.MustCompile("^[\u30A1-\u30FC]{5}$")
	game := &Game{TurnsRemaining: 5}
	set_answer(game)
	game.print_empty()
	//入力を受け付ける
MainLoop:
	for {
		guess := &Guess{Background: make([]string, N)}
		var poke string
		fmt.Print("予想->")
		fmt.Scanf("%s", &poke)
		if !valid.Match([]byte(poke)) {
			fmt.Print("カタカナ5文字を入力してください")
			fmt.Printf("\033[1F\033[0K")
			continue MainLoop
		}
		guess.Pokemon = []rune(poke)
		game.TurnsRemaining--

		game.comparing(guess)
		game.display(guess)

		if string(guess.Pokemon) == string(game.Answer) {
			game.Won = true
			break MainLoop
		}
		if game.TurnsRemaining == 0 {
			break MainLoop
		}
	}

	if game.Won {
		fmt.Printf("クリア成功! 正解は%s!\n", string(game.Answer))
	} else {
		fmt.Printf("クリア失敗! 正解は%s!\n", string(game.Answer))
	}

}
