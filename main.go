package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

func main() {
	// os.Setenv("HTTP_PROXY", "10.4.200.197:8089")

	// get categories from json file
	file, _ := ioutil.ReadFile("categories.json")

	data := Categories{}

	_ = json.Unmarshal([]byte(file), &data)

	crawlAllFromCategories(data)
}

func crawlAllFromCategories(categories Categories) {
	var wg sync.WaitGroup

	jobs := make(chan Category, 100)

	wg.Add(1)
	go worker(jobs, &wg)

	// schedule to run program each 3 hour
	for true {
		for i := 0; i < len(categories.List); i++ {
			jobs <- categories.List[i]
		}
		time.Sleep(5 * time.Hour)
	}

	close(jobs)

	wg.Wait()
}

func worker(jobs <-chan Category, wg *sync.WaitGroup) {
	defer wg.Done()

	// create output directory
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.Mkdir("./output", 0755)
	}

	for j := range jobs {
		// open or create file
		dt := time.Now()
		f, _ := os.OpenFile("./output/"+j.Title+"___"+dt.Format("20060102150405")+".json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		fmt.Println("processing job: ", j)
		crawlFromCategory(j, f)
	}
}

func crawlFromCategory(category Category, f *os.File) {
	// open leveldb connection
	db := createOrOpenDb("./db/" + category.Title)
	defer db.Close()

	users := NewUsers()
	res := getHTMLPage(category.URL)

	//handle error
	if res == nil {
		return
	}

	users.getAllUserInformation(res, category.Title, f, db)
	users.TotalPages++

	prevPage := category.URL

	for i := 2; i <= 30; i++ {
		users.TotalPages++
		nextPageLink := users.getNexURL(res)

		if prevPage == nextPageLink {
			println("End of Category: " + category.Title)
			break
		} else {
			prevPage = nextPageLink
		}

		if nextPageLink == "" {
			break
		}

		res = getHTMLPage(nextPageLink)

		//handle error
		if res == nil {
			break
		}

		users.getAllUserInformation(res, category.Title, f, db)
	}
}
