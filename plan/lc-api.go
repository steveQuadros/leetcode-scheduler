package plan

import (
    "net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"github.com/valyala/fastjson"
)

const url = "https://leetcode.com/graphql"

func Fetch(questionNames []string) ([]*Question, error) {
	var fetched []*Question

    for _, question := range questionNames {
        query := getQuery(question)
		queryString, err := json.Marshal(query)
		if err != nil {
			return fetched, err
		}

		queryBytes := []byte(queryString)
		rw := ioutil.NopCloser(bytes.NewBuffer(queryBytes))
        req, err := http.NewRequest(http.MethodGet, url, rw)
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			return fetched, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fetched, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fetched, err
		}

        if resp.StatusCode != 200 {
            return fetched, err
		}
		
		title := fastjson.GetString(body,"data", "question", "titleSlug")
		link := "https://leetcode.com/problems/" + fastjson.GetString(body,"data", "question", "titleSlug")
		difficulty := fastjson.GetString(body,"data", "question", "difficulty")
		fetched = append(fetched, &Question{Title: title, Link: link, Difficulty: difficulty})
	}
    
	return fetched, nil
}

func getQuery(question string) map[string]interface{} {
	return map[string]interface{}{
        "operationName": "questionData",
        "variables": map[string]interface{}{
            "titleSlug": question,
        },
        "query": "query questionData($titleSlug: String!) {question(titleSlug: $titleSlug) {    questionId    questionFrontendId   title    titleSlug   isPaidOnly    difficulty  }}",
    }
}