package plan

// func Update(plans []*Plan, intervals []int) []*Plan {
// 	schedule := map[string]*Plan{}
// 	for _, p := range plans {
// 		plansMap[plans.Date] = p
// 	}

// 	for _, p := range schedule {
// 		for _, q := range p.New {
// 			for _, i := range intervals {
// 				tempDate := curDate.AddDate(0,0,i-1)
// 				dateKey := formatTime(tempDate)
// 				_, ok := schedule[dateKey]
// 				if !ok {
// 					schedule[dateKey] = &Plan{Date: dateKey}
// 				}
// 				questions := schedule[dateKey]
// 				questions.Review = append(questions.Review, q)
// 			}
// 		}
// 	}
// }