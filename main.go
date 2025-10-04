package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Questions struct {
	Options []string
	Answer  int
	Text    string
}

type GameState struct {
	Nome     string
	Temas    int
	Points   int
	Question []Questions
}

func (g *GameState) Init() {
	fmt.Println("Bem vindo ao jogo!\n Você terá 10 segundos para responder as perguntas.")
	fmt.Println("Qual é o seu nome?")

	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')

	if err != nil {
		panic("Erro ao ler a string")
	}

	g.Nome = name

	fmt.Printf("bem vindo ao jogo, %s\n", g.Nome)

	temas := []string{"História", "Matemática", "Geografia"}

	fmt.Println("Qual tema irá escolher?")
	for index := range temas {
		fmt.Printf("[%d] %s\n", index+1, temas[index])
	}

	temaEscolhido, _ := reader.ReadString('\n')

	g.Temas, _ = toInt(temaEscolhido[:len(temaEscolhido)-2])
}

func (g *GameState) ProccessCSV() {
	var f *os.File
	var err error
	var fileName string

	switch g.Temas - 1 {
	case 0:
		fileName = "historia.csv"
	case 1:
		fileName = "matematica.csv"
	case 2:
		fileName = "geografia.csv"
	default:
		panic("tema inválido!")
	}

	f, err = os.Open(filepath.Join("Temas", fileName))

	defer f.Close()

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()

	if err != nil {
		panic(err.Error())
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Questions{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}
			g.Question = append(g.Question, question)
		}
	}
}
func main() {
	game := &GameState{Points: 0}
	game.Init()
	game.ProccessCSV()

	game.Run()

}

func (g *GameState) Run() {
	// Exibir pergunta
	for index, question := range g.Question {
		fmt.Printf("%d) %s\n", index+1, question.Text)
		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", j+1, option)
		}
		fmt.Println("Digite uma alternativa: ")

		var answer int
		var err error

		read := make(chan string)

		for {
			go func() {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				read <- text
			}()

			select {
			case input := <-read:
				answer, err = toInt(input[:len(input)-2])
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
			case <-time.After(10 * time.Second):
				fmt.Println("Tempo de 10 segundos encerrado, a questão será marcada como errada.")
			}

			if answer == question.Answer {
				g.Points += 10
			}
			break
		}

	}

	fmt.Printf("a quantidade total de pontos foi: %d", g.Points)

	if g.Points >= 60 {
		fmt.Println("\nVocê passou no teste, parabéns!")
	} else if g.Points < 60 {
		fmt.Println("\nVocê não passou no teste, tente novamente!")
	}
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("não é permitido caractere diferente de número, retorne um número")

	}

	return i, nil

}
