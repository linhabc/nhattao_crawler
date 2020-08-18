package main

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (categories *Categories) getAllCategories(doc *goquery.Document) {

	doc.Find(".listCatgoryMenu-item").Each(func(i int, s *goquery.Selection) {
		data, _ := s.Attr("data-children")
		if len(data) == 2 {
			return
		}
		dec := json.NewDecoder(strings.NewReader(data))
		var cat TmpCategories

		if err := dec.Decode(&cat); err == io.EOF {
			return
		} else if err != nil {
			log.Fatal(err)
		}

		for _, v := range cat {
			category := Category{
				Title: v.Title,
				URL:   "https://nhattao.com/" + v.URL,
			}

			categories.Total++
			categories.List = append(categories.List, category)
		}

	})
}

// run this function to crawlAllCategories
// func main() {
// 	categories := newCategories()
// 	res := getHTMLPage("https://nhattao.com/")

// 	categories.getAllCategories(res)

// 	userJSON, err := json.Marshal(categories)
// 	checkError(err)
// 	err = ioutil.WriteFile("./categories.json", userJSON, 0644) // Ghi dữ liệu vào file JSON
// 	checkError(err)
// }
