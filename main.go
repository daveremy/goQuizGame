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
)

type problem struct {
	Text   string `json:"text"`
	Answer string `json:"answer"`
}

func main() {
	// Handle command line options
	problemsFileFlag := flag.String("problemsFile", "problems.csv", "Problems File Name")
	flag.Parse()

	problems := readProblemsFromCsv(*problemsFileFlag)

	var numCorrect, numIncorrect int = 0, 0
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Quiz Master")
	fmt.Printf("There are %d problems in the quiz.  Ready? (Y/N) ", len(problems))
	response, _ := reader.ReadString('\n')
	if start := strings.TrimSpace(response); strings.EqualFold("Y", start) {
		fmt.Println("Ok, let's do this!")
		for i, problem := range problems {
			fmt.Println("----------------------------------------------------")
			fmt.Printf("%d: %s\n", i, problem.Text)
			response, _ := reader.ReadString('\n')
			if answer := strings.TrimSpace(response); strings.EqualFold(answer, problem.Answer) {
				numCorrect++
				fmt.Println("Correct")
			} else {
				numIncorrect++
				fmt.Println("Incorrect")
			}
		}
		fmt.Println("----------------------------------------------------")
		fmt.Println("Finished!  Summary:")
		fmt.Printf("You got %d Correct and %d Incorrect (%v%%)\n", numCorrect, numIncorrect, math.Round(float64(numCorrect)/float64(len(problems))*100))
		fmt.Println("----------------------------------------------------")
	} else {
		fmt.Printf("Ok, maybe later! (You typed '%s' rather than 'Y')", start)
	}
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
		// fmt.Println(record)
		problems = append(problems, problem{
			Text:   record[0],
			Answer: record[1],
		})
		// problemsJSON, _ := json.Marshal(problems)
		// fmt.Println(string(problemsJSON))
	}
	return problems
}
