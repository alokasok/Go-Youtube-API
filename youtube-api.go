package main

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	URL     = "https://www.googleapis.com/youtube/v3/"
	API_KEY = "AIzaSyBi1ieWcovv1p28szRaj1f_hLC2nAp5dkg"
)

var url string

var info map[string]string

func main() {
	fmt.Println("Youtube API in go")

	flag.StringVar(&url, "url", "test", "url to get the info")
	flag.Parse()
	fmt.Println("URL is " + url)
	if url == "test" {
		log.Fatal("Youtube URL not specified")
	}
	r, _ := regexp.Compile("v=")

	fmt.Println(r.MatchString(url))
	substring := r.FindStringIndex(url)
	length := len(url)
	fmt.Println(length)
	slice := url[substring[1]:length]
	fmt.Println(slice)

	if _, err := os.Stat("youtube.db"); err == nil {
		fmt.Println("DB exists")
	} else {
		fmt.Println("DB doesnt exist creating a new db")
		c, _ := sqlite3.Open("youtube.db")
		c.Exec("CREATE TABLE video_info(url,title,info,likecount,dislikecount,published_date)")
	}

	resp, err := http.Get(URL + "videos?id=" + slice + "&key=" + API_KEY + "&part=snippet,contentDetails,statistics,status")
	if err != nil {
		log.Fatal("API failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	var f interface{}
	error := json.Unmarshal(body, &f)
	if error != nil {
		log.Fatal("JSON parsing failed")
	}
	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			//fmt.Println(k, "is an array:")
			info := vv
			//fmt.Println(info)
			for i, u := range vv {
				fmt.Println("inside interface")
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}

	for k, v := range info {
		fmt.Println(k)
		fmt.Println(v)
	}

	/*for k := range m {
		fmt.Println(k)
	}*/

}
