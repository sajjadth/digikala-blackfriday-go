package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func errorHandeler(code int, err error) {
	if err != nil {
		log.Fatal(code, err)
	}
}

func getAllLinks() []string {
	var links []string
	baseUrl := "https://www.digikala.com/treasure-hunt/products/??has_selling_stock=1&pageno=%d&sortby=4"
	from, to := 1, 42

	for from <= to {
		currentUrl := fmt.Sprintf(baseUrl, from)
		res, err := http.Get(currentUrl)
		errorHandeler(1, err)
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		errorHandeler(2, err)
		doc.Find(".c-promotion-box__image").Each(func(i int, s *goquery.Selection) {
			link, haveLink := s.Attr("href")
			if haveLink {
				newLink := fmt.Sprintf("https://www.digikala.com%v", link)
				links = append(links, newLink)
			}
		})
		from++
	}
	return links
}

func getAllImages(links []string) []string {
	var images []string
	for _, link := range links {
		res, err := http.Get(link)
		errorHandeler(3, err)
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		errorHandeler(4, err)
		doc.Find(".pannable-image").Each(func(i int, s *goquery.Selection) {
			img, haveSrc := s.Attr("data-src")
			if haveSrc {
				imgDate, _ := strconv.Atoi(string(img[97:107]))
				if imgDate > 1637800000 {
					fmt.Println(img)
					images = append(images, img)
				}
			}
		})
	}
	return images
}

func treasureHunt(images []string) []string {
	treasure := make([]string, 15)
	for _, img := range images {
		imgDate, _ := strconv.Atoi(string(img[97:107]))
		for i, v := range treasure {
			if len(v) == 0 {
				treasure[i] = img
			} else {
				treasureDate, _ := strconv.Atoi(string(v[97:107]))
				if treasureDate < imgDate {
					treasure[i] = img
				}
			}
		}
	}

	return treasure
}

func main() {
	start := time.Now()
	links := getAllLinks()
	images := getAllImages(links)
	treasure := treasureHunt(images)
	for _, t := range treasure {
		err := exec.Command("rundll32", "url.dll,FileProtocolHandler", t).Start()
		errorHandeler(5, err)
	}
	fmt.Println(time.Since(start))
}
