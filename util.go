package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/syndtr/goleveldb/leveldb"
)

func getHTMLPage(url string) *goquery.Document {
	// Request the HTML page.
	// res, err := http.Get(url)

	res, err := http.Get(url)

	if err != nil {
		println("ERROR GET")
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode == 429 {
		for {
			time.Sleep(1 * time.Second)
			res, err = http.Get(url)
			if res.StatusCode == 200 {
				break
			}
		}
	}

	if res.StatusCode != 200 {
		print("ERORR RES STATUS: ")
		return nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil
	}
	return doc
}

func (users *Users) getNexURL(doc *goquery.Document) string {
	aLink := doc.Find("a.text:last-child")
	nextPageLink, _ := aLink.Attr("href")

	// Trường hợp không có url
	if nextPageLink == "" {
		println("End of Category")
		return ""
	}

	nextPageLink = rootLink + nextPageLink

	print("NEXTPAGE: ")
	println(nextPageLink)

	time.Sleep(1 * time.Second)
	return nextPageLink
}

func (users *Users) getAllUserInformation(doc *goquery.Document, category string, f *os.File, db *leveldb.DB) {
	doc.Find(".Nhattao-CardItem--inner a.title").Each(func(i int, s *goquery.Selection) {
		userLink, _ := s.Attr("href")
		userLink = rootLink + userLink
		users.getUserInformation(userLink, category, f, db)
	})
}

func (users *Users) getUserInformation(url string, category string, f *os.File, db *leveldb.DB) {

	res := getHTMLPage(url)
	if res == nil {
		return
	}

	currentTime := time.Now()
	time := currentTime.Format("01/02/2006")
	userName := res.Find(".threadview-header--seller a span").Text()
	title := res.Find("h2.threadview-header--title").Text()
	price := res.Find(".threadview-header--classifiedPrice").Text()
	location := res.Find("dd span.address").Text()
	phoneNum := res.Find("#nhattao2019-contactPhone").Text()

	userName = strings.TrimSpace(userName)
	phoneNum = strings.TrimSpace(phoneNum)
	title = strings.TrimSpace(title)
	time = strings.TrimSpace(time)
	location = strings.TrimSpace(location)
	price = strings.TrimSpace(price)

	time = strings.Replace(time, "/", "-", 2)

	if len(phoneNum) == 0 {
		println("phone num = 0 " + url)
		return
	}

	splitResult := strings.Split(url, ".")
	tmpid := splitResult[len(splitResult)-1]

	splitResult = strings.Split(tmpid, "/")
	id := splitResult[0]

	// check if id is exist in db or not
	checkExist := getData(db, id)
	if len(checkExist) != 0 {
		println("Exist: " + id)
		return
	}
	println("None_exist: " + id)

	user := User{
		ID:          id,
		PhoneNumber: phoneNum,
		UserName:    userName,
		Title:       title,
		Time:        time,
		Location:    location,
		Price:       price,
	}

	_ = putData(db, id, phoneNum)

	// convert User to JSON
	userJSON, err := json.Marshal(user)

	checkError(err)
	io.WriteString(f, string(userJSON)+"\n")

	users.TotalUsers++
	users.List = append(users.List, user)
}

func checkError(err error) {
	if err != nil {
		print("Error: ")
		log.Println(err)
	}
}

const rootLink = "https://nhattao.com/"
