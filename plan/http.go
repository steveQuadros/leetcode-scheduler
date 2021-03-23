package plan

import (
    "log"
    "net/http"
	"html/template"
)

func handler(w http.ResponseWriter, r *http.Request) {
    
}

func Serve(sp *StudyPlan) {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("plan.html"))
		t.Execute(w, sp)
	})
    log.Fatal(http.ListenAndServe(":8080", nil))
}