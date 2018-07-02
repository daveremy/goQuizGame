package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type problem struct {
	Text   string `json:"text"`
	Answer string `json:"answer"`
}

func main() {
	// Handle command line options
	problemsFileFlag := flag.String("problemsFile", "problems.csv", "Problems File Name")
	timeLimitFlag := flag.Int("timeLimit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	problems := readProblemsFromCsv(*problemsFileFlag)
	doQuiz(problems, *timeLimitFlag)
}

func readProblemsFromCsv(filename string) []problem {
	csvFile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var problems []problem
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		problems = append(problems, problem{
			Text:   record[0],
			Answer: record[1],
		})
	}
	return problems
}

func doQuiz(problems []problem, timeLimit int) {
	var numCorrect = 0
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Quiz Master")
	fmt.Printf("There are %d problems in the quiz, you have %v seconds to complete it.  Ready? (Y/N)\n", len(problems), timeLimit)
	timer := time.NewTimer(time.Second * time.Duration(timeLimit))
	response, _ := reader.ReadString('\n')
	if start := strings.TrimSpace(response); strings.EqualFold("Y", start) {
		fmt.Println("Ok, let's do this!")
	problemLoop:
		for i, problem := range problems {
			fmt.Println("----------------------------------------------------")
			fmt.Printf("%d: %s\n", i+1, problem.Text)
			responseCh := make(chan string)
			go func() {
				var response string
				fmt.Scanf("%s\n", &response)
				responseCh <- strings.TrimSpace(response)
			}()
			select {
			case <-timer.C:
				fmt.Println("----------------------------------------------------")
				fmt.Printf("\nQuiz timed out after %v seconds.", timeLimit)
				break problemLoop
			case answer := <-responseCh:
				if strings.EqualFold(answer, problem.Answer) {
					numCorrect++
					fmt.Println("Correct")
				} else {
					fmt.Printf("Incorrect, should be: %s\n", problem.Answer)
				}
			}
		}
		fmt.Println("----------------------------------------------------")
		fmt.Println("Finished!  Summary:")
		fmt.Printf("You got %d correct and %d incorrect (or unanswered) (%v%%)\n", numCorrect, len(problems)-numCorrect, math.Round(float64(numCorrect)/float64(len(problems))*100))
	} else {
		fmt.Printf("Ok, maybe later! (You typed '%s' rather than 'Y')", start)
	}
}
