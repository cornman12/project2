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

	manga := Manga{"anh-hung-ta-khong-lam-lau-roi","http://www.nettruyen.com/truyen-tranh/anh-hung-ta-khong-lam-lau-roi"}
	
	
	go getChapterList(manga, chapterchan)


	for chap := range chapterchan {

		fmt.Println(chap)
		//go getChapter(chap)
	}
	

	fmt.Scanln()

	//getChapter("http://www.nettruyen.com/truyen-tranh/anh-hung-ta-khong-lam-lau-roi/chap-1/409746")



}


type Manga struct{
	name string
	url string
}



type Chapter struct {
	mangaName string
	name string
	url  string
}



func getChapter(mangaChapter Chapter ) {

	mangaName := mangaChapter.mangaName
	chapterName := mangaChapter.name
	url := mangaChapter.url
	pwd, _ := os.Getwd()

	fmt.Println("Getting images from "+ url)
	httpconfig := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	httpclient := &http.Client{Transport: httpconfig}
	resp, err := httpclient.Get(url)
	defer resp.Body.Close()
	if err !=nil{
		return
	}
	htmltoken := html.NewTokenizer(resp.Body)
	for{

		tokenTag := htmltoken.Next()
		if tokenTag != html.ErrorToken{
			nextToken := htmltoken.Token()
			if nextToken.Data == "img"{
				
				for _, att := range nextToken.Attr{
					if att.Key =="data-original"{
						fmt.Println(nextToken.Attr[1].Val)
						fmt.Println(nextToken.Attr[2].Val)
					}
				}
			}

		}else{
			return
		}


	}

}








func getChapterList(manga Manga, chapterchan chan Chapter) {
	mangaName:= manga.name
	url := manga.url
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
				chapter := Chapter{mangaName ,t.Data, chapterurl}
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
