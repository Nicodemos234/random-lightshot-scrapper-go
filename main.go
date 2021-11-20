package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getRandomLink() string {
	baseUrl := "https://prnt.sc/"
	return baseUrl + randSeq(6)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	client := &http.Client{}
	i := 0
	for {
		newLink := getRandomLink()
		req, err := http.NewRequest("GET", newLink, nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		bodyString := string(body)
		r, _ := regexp.Compile("<img class=\"no-click screenshot-image\" src=\"(.*)+\" crossorigin=\"anonymous\"")
		match := r.FindStringSubmatch(bodyString)

		if match != nil {
			if !strings.Contains(match[0], "//st.prntscr.com/2021/10/22/2139/img/0_173a7b_211be8ff.png") {
				clearedString := strings.Replace(match[0], "<img class=\"no-click screenshot-image\" src=\"", "", -1)
				clearedString = strings.Replace(clearedString, "\" crossorigin=\"anonymous\"", "", -1)
				fmt.Print(clearedString, ", ")
				DownloadFile("img/"+strconv.Itoa(i)+".png", clearedString)
				i++
			}
		}
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
