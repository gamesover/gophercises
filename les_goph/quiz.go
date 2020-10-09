package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type question struct {
	questionText string
	answer       string
}

func main() {
	csvPtr := flag.String("csv", "les_goph/problems.csv", "csv file in format: question,answer")
	timeLimitPtr := flag.Int("limit", 30, "the time limit in seconds")
	flag.Parse()

	csvfile, err := os.Open(*csvPtr)
	if err != nil {
		exit(fmt.Sprintf("Couldn't open the csv file %s", *csvPtr))
	}

	problems := parseLines(csvfile)
	correctCount := 0

	timer := time.NewTimer(time.Duration(*timeLimitPtr) * time.Second)
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.questionText)
		answerCh := make(chan string)
		go func() {
			var input string
			fmt.Scanln(&input)
			answerCh <- input
		}()
		select {
		case <-timer.C:
			fmt.Printf("\n\nYou scored %d out of %d\n", correctCount, len(problems))
			return
		case input := <-answerCh:
			answer := strings.TrimSpace(input)
			if answer == problem.answer {
				correctCount++
			}
		}

	}

	fmt.Printf("\n\nYou scored %d out of %d\n", correctCount, len(problems))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(file *os.File) []question {
	var questions []question
	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			exit("fail to read file")
		}
		newQuestion := question{record[0], strings.TrimSpace(record[1])}
		questions = append(questions, newQuestion)
	}

	return questions
}
