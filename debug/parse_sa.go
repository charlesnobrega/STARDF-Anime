package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	f, err := os.Open("superanimes_results_v2.html")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		fmt.Println("Error parsing:", err)
		return
	}

	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		action, _ := s.Attr("action")
		method, _ := s.Attr("method")
		id, _ := s.Attr("id")
		
		fmt.Printf("Form id='%s', action='%s', method='%s'\n", id, action, method)
		s.Find("input").Each(func(j int, input *goquery.Selection) {
			name, _ := input.Attr("name")
			inputType, _ := input.Attr("type")
			fmt.Printf("  Input name='%s', type='%s'\n", name, inputType)
		})
	})
}
