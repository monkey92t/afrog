package main

import (
	"log"

	"github.com/monkey92t/afrog/v2/pkg/report"
	"github.com/monkey92t/afrog/v2/pkg/result"
)

func main() {
	filename := "xxx.htm"
	report, err := report.NewReport(filename, report.DefaultTemplate)
	if err != nil {
		log.Fatalf("newReprot err: %v", err)
	}
	report.Result = &result.Result{IsVul: true, Target: "http://localhost"}
	err = report.Append("1")
	if err != nil {
		log.Fatalf("Append err: %v", err)
	}
}
