package plan

import (
	"time"
	"sort"
)

type Question struct {
	Link, Title, Difficulty string
	InitialDate time.Time
	ReviewDates []time.Time
}

type Plan struct {
	Date string
	Completed bool
	Questions
}

type Questions struct {
	New []*Question
	Review []*Question
}

var Difficulty = map[string]int{
	"Easy": 1,
	"Medium": 2,
	"Hard": 3,
}

func GenerateDates(questions []*Question, intervals []int, startDate time.Time, questionsPerDay int) []*Question {
	// sort by difficult
	sort.Slice(questions, func(i,j int) bool { 
		return Difficulty[questions[i].Difficulty] < Difficulty[questions[j].Difficulty]
	})

	var count int
	var curDate = startDate
	for _, question := range questions {
		if count >= questionsPerDay {
			curDate = curDate.AddDate(0,0,1)
			count = 0
		}
		count++
		question.InitialDate = curDate

		for _, interval := range intervals {
			question.ReviewDates = append(question.ReviewDates, curDate.AddDate(0, 0, interval-1))
		}
	}
	return questions	
}

func ByDate(questions []*Question) []*Plan {
	var schedule = map[string]*Plan{}
	for _, q := range questions {
		dateKey := formatTime(q.InitialDate)
		_, ok := schedule[dateKey]
		if !ok {
			schedule[dateKey] = &Plan{Date: dateKey}
		}
		plan := schedule[dateKey]
		plan.New = append(plan.New, q)

		for _, d := range q.ReviewDates {
			dateKey := formatTime(d)
			_, ok := schedule[dateKey]
			if !ok {
				schedule[dateKey] = &Plan{Date: dateKey}
			}
			plan := schedule[dateKey]
			plan.Review = append(plan.Review, q)
		}
	}
	
	dates := []string{}
	for date, _ := range schedule {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	plans := []*Plan{}
	for _, d := range dates {
		plans = append(plans, schedule[d])
	}

	return plans
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02")
}

// func print() {
// 	details := &strings.Builder{}
// 	overview := &strings.Builder{}
// 	fmt.Fprintln(overview, "Overview:")
// 	var max int
// 	var maxDate string
// 	for _, date := range dates {
// 		questions := schedule[date]
// 		new := questions.New
// 		review := questions.Review
// 		fmt.Fprintln(details, date)
// 		curTotal := len(new) + len(review)

// 		fmt.Fprintln(details, "Total Problems", curTotal)
// 		fmt.Fprintf(overview, "Date: %s, NumQuestions: %d\n", date, curTotal)

// 		if curTotal > max {
// 			max = curTotal
// 			maxDate = date
// 		}

// 		fmt.Fprintln(details, "New: ", len(new))
// 		for _, q := range new {
// 			fmt.Fprintln(details, q.Title, " - ", q.Difficulty, " - ", q.Link)
// 		}

// 		fmt.Fprintln(details, "Review: ", len(review))
// 		for _, q := range review {
// 			fmt.Fprintln(details, q.Title, " - ", q.Difficulty, " - ", q.Link)
// 		}
// 		fmt.Fprintln(details)
// 	}

// 	header := &strings.Builder{}
// 	fmt.Fprintf(header, "Start: %s\n", formatTime(startDate))
// 	fmt.Fprintf(header, "End: %s\n", dates[len(dates)-1])
// 	fmt.Fprint(header, "Intervals: ") 
// 	for _, i := range intervals {
// 		fmt.Fprintf(header, " %d", i)
// 	}
// 	fmt.Fprintln(header)
// 	fmt.Fprintf(header, "Questions Per Day: %d\n", questionsPerDay)
// 	fmt.Fprintf(header, "Max questions: %d on %s\n", max, maxDate)

// 	fmt.Println(header.String())
// 	fmt.Println(overview.String())
// 	fmt.Println(details.String())
// }