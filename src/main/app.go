package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
	//"net/http"
)

func main() {
	chapterchan := make(chan Chapter)
	fmt.Println("hello world")
	go getRequest("http://www.nettruyen.com/truyen-tranh/anh-hung-ta-khong-lam-lau-roi", chapterchan)

	for chap := range chapterchan {

		fmt.Println(chap)
	}

}

type Chapter struct {
	name string
	url  string
}

func getRequest(url string, chapterchan chan Chapter) {
	chapterURLChannel := make(chan string, 1)
	httpconfig := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	httpclient := &http.Client{Transport: httpconfig}
	fmt.Println("getting data from :" + url)
	chapterelement := 0
	resp, _ := httpclient.Get(url)
	defer resp.Body.Close()
	htmltoken := html.NewTokenizer(resp.Body)

	for {
		tt := htmltoken.Next()
		//fmt.Println(tt)
		if tt != html.ErrorToken {

			t := htmltoken.Token()
			if tt.String() == "Text" && chapterelement == 1 {
				chapterelement = 0
				chapterurl := <-chapterURLChannel
				chapter := Chapter{t.Data, chapterurl}
				chapterchan <- chapter
				//fmt.Println(t.Data)
				//fmt.Println(chapterurl)
			}
			if t.Data == "a" {
				for _, att := range t.Attr {
					if att.Key == "data-id" {
						chapterelement = 1
						chapterURLChannel <- t.Attr[0].Val
						//fmt.Println(t.Attr[0].Val)

					}
				}

			} else {
				chapterelement = 0
			}

		} else {
			close(chapterURLChannel)
			fmt.Println("ended")
			close(chapterchan)
			return
		}
	}

}

func saveImage(url string) {
	httpconfig := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	httpclient := &http.Client{Transport: httpconfig}
	fmt.Println("getting data from :" + url)
	resp, err := httpclient.Get(url)
	defer resp.Body.Close()
	out, _ := os.Create("test.jpg")
	_, errx := io.Copy(out, resp.Body)
	if errx != nil {
		fmt.Println("error")
	}
	if err != nil {
		fmt.Println(err.Error())

	}
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		fmt.Println(string(body))
	}

}
