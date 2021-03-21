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
var Difficulty = map[string]int{
	"Easy": 1,
	"Medium": 2,
	"Hard": 3,
}

type Question struct {
	Link, Title, Difficulty string
}

type Questions struct {
	New []Question
	Review []Question
}

func main() {
	c := csv.NewReader(os.Stdin)
	// skip header
	_, err := c.Read()
	if err != nil {
		panic(err)
	}
	
	var allQuestions []Question
	for line, err := c.Read(); err != io.EOF; line, err = c.Read() {
		if err != nil {
			panic(err)
		}
		allQuestions = append(allQuestions, Question{Title: line[0], Link: line[1], Difficulty: line[2]})
	}

	if len(allQuestions) < 1 {
		fmt.Println("no questions to schedule")
		return
	}

	// sort by difficult
	sort.Slice(allQuestions, func(i,j int) bool { 
		return Difficulty[allQuestions[i].Difficulty] < Difficulty[allQuestions[j].Difficulty]
	})

	var schedule = map[string]*Questions{}
	var curQuestions = 0
	var curDate = startDate
	for _, q := range allQuestions {
		if curQuestions > questionsPerDay-1 {
			curDate = curDate.AddDate(0,0,1)
			curQuestions = 0
		}
		curQuestions++

		dateKey := formatTime(curDate)
		_, ok := schedule[dateKey]
		if !ok {
			schedule[dateKey] = &Questions{}
		}
		questions := schedule[dateKey]
		questions.New = append(questions.New, q)

		// for each interval add interval to innerDate and add question to map[innerdate]
		for _, i := range intervals {
			tempDate := curDate.AddDate(0,0,i-1)
			dateKey := formatTime(tempDate)
			_, ok := schedule[dateKey]
			if !ok {
				schedule[dateKey] = &Questions{}
			}
			questions := schedule[dateKey]
			questions.Review = append(questions.Review, q)
		}
	}
	
	dates := []string{}
	for date, _ := range schedule {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	details := &strings.Builder{}
	overview := &strings.Builder{}
	fmt.Fprintln(overview, "Overview:")
	var max int
	var maxDate string
	for _, date := range dates {
		questions := schedule[date]
		new := questions.New
		review := questions.Review
		fmt.Fprintln(details, date)
		curTotal := len(new) + len(review)

		fmt.Fprintln(details, "Total Problems", curTotal)
		fmt.Fprintf(overview, "Date: %s, NumQuestions: %d\n", date, curTotal)

		if curTotal > max {
			max = curTotal
			maxDate = date
		}

		fmt.Fprintln(details, "New: ", len(new))
		for _, q := range new {
			fmt.Fprintln(details, q.Title, " - ", q.Difficulty)
		}

		fmt.Fprintln(details, "Review: ", len(review))
		for _, q := range review {
			fmt.Fprintln(details, q.Title, " - ", q.Difficulty)
		}
		fmt.Fprintln(details)
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
	fmt.Println(overview.String())
	fmt.Println(details.String())
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02")
}