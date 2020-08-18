package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/syndtr/goleveldb/leveldb"
)

func getHTMLPage(url string) *goquery.Document {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	req.Host = "nhattao.com"
	req.Header = map[string][]string{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:77.0) Gecko/20100101 Firefox/77.0"},
		// "Accept-Encoding": {"gzip, deflate"},

		// "Accept-Language": {"vi-VN", "vi", "q=0.8,en-US", "q=0.5,en", "q=0.3"},
		"Cookie": {"nhattao_session=05b69b4f16ec28215ef29598a9996f60; xf_vim|mudim-settings=26; G_ENABLED_IDPS=google; _cfduid=d087a49f7192613603e541a3d0337f0e91597744875; ga=GA1.2.1555371478.1597744810; gid=GA1.2.1217142350.1597744810; gat=1; _gads=ID=50fb02670837f4f0-22570feffec200ec:T=1597744877:S=ALNI_MYu6AUNqV62AHSYzaErs_rUmWkDxQ; fbp=fb.1.1597744810408.1159376184"},

		"Referer": {"https://nhattao.com/"},
	}
	client := &http.Client{}

	// Request the HTML page.
	// res, err := http.Get(url)

	res, err := client.Do(req)
	if err != nil {
		println("ERROR GET")
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		print("ERORR RES STATUS: ")
		println(res.StatusCode)
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
	aLink := doc.Find("a.text")
	nextPageLink, _ := aLink.Attr("href")

	// Trường hợp không có url
	if nextPageLink == "" {
		println("End of Category")
		return ""
	}

	nextPageLink = rootLink + nextPageLink

	print("NEXTPAGE: ")
	println(nextPageLink)

	time.Sleep(5 * time.Second)
	return nextPageLink
}

func (users *Users) getAllUserInformation(doc *goquery.Document, category string, f *os.File, db *leveldb.DB) {
	var wg sync.WaitGroup
	doc.Find(".Nhattao-CardItem--inner a.title").Each(func(i int, s *goquery.Selection) {
		userLink, _ := s.Attr("href")
		wg.Add(1)
		userLink = rootLink + userLink
		go users.getUserInformation(userLink, category, &wg, f, db)
	})
	wg.Wait()
}

func (users *Users) getUserInformation(url string, category string, wg *sync.WaitGroup, f *os.File, db *leveldb.DB) {
	defer wg.Done()

	time.Sleep(5 * time.Second)
	res := getHTMLPage(url)
	if res == nil {
		return
	}

	currentTime := time.Now()
	time := currentTime.Format("09-07-2017")
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
