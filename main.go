package main

import (
	"encoding/csv"
	"os"
	"io"
	"time"
	"fmt"
	"sort"
	"strings"
	// "flags"
)

var intervals = []int{2,4,8,16}
var startDate = time.Now()
var questionsPerDay = 3

type Question struct {
	Link, Title, Difficulty string
}

type Questions struct {
	New []Question
	Review []Question
}

func main() {
	c := csv.NewReader(os.Stdin)
	schedule := map[string]Questions{}

	var curQuestions = 0
	var curDate = startDate
	// skip header
	_, err := c.Read()
	if err != nil {
		panic(err)
	}
	for line, err := c.Read(); err != io.EOF; line, err = c.Read() {
		if err != nil {
			panic(err)
		}

		// go to next day if we've reached max new questions for this day
		if curQuestions > questionsPerDay-1 {
			curDate = curDate.AddDate(0,0,1)
			curQuestions = 0
		}
		curQuestions++
		question := Question{Title: line[0], Link: line[1], Difficulty: line[2]}
		questions := schedule[formatTime(curDate)] 
		questions.New = append(questions.New, question)
		schedule[formatTime(curDate)] = questions

		// for each interval add interval to innerDate and add question to map[innerdate]
		for _, i := range intervals {
			tempDate := curDate.AddDate(0,0,i-1)
			reviews := schedule[formatTime(tempDate)]
			reviews.Review = append(reviews.Review, question)
			schedule[formatTime(tempDate)] = reviews 
		}
	}
	
	dates := []string{}
	for date, _ := range schedule {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	b := &strings.Builder{}
	var max int
	var maxDate string
	for _, date := range dates {
		questions := schedule[date]
		new := questions.New
		review := questions.Review
		fmt.Fprintln(b, date)
		curTotal := len(new) + len(review)
		fmt.Fprintln(b, "Total Problems", curTotal)
		if curTotal > max {
			max = curTotal
			maxDate = date
		}

		fmt.Fprintln(b, "New: ", len(new))
		for _, q := range new {
			fmt.Fprintln(b, q.Title, " - ", q.Difficulty)
		}

		fmt.Fprintln(b, "Review: ", len(review))
		for _, q := range review {
			fmt.Fprintln(b, q.Title, " - ", q.Difficulty)
		}
		fmt.Fprintln(b)
	}

	header := &strings.Builder{}
	fmt.Fprintf(header, "Start: %s\n", formatTime(startDate))
	fmt.Fprintf(header, "End: %s\n", dates[len(dates)-1])
	fmt.Fprint(header, "Intervals: ") 
	for _, i := range intervals {
		fmt.Fprintf(header, " %d", i)
	}
	fmt.Fprintln(header)
	fmt.Fprintf(header, "Questions Per Day: %d\n", questionsPerDay)
	fmt.Fprintf(header, "Max questions: %d on %s\n", max, maxDate)
	fmt.Println(header.String())

	fmt.Println(b.String())
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02")
}