package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var netClient = &http.Client{
	Timeout: 60 * time.Second,
}

func main() {
	http.HandleFunc("/", serveTemplate)

	fmt.Println("Listening on :80...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The host is ", r.Host)

	// Choose witch port to fetch from
	var port string
	switch r.Host {
	case "ebinbellini.top":
		port = "9001"
	case "www.ebinbellini.top":
		port = "9001"
	case "chat.ebinbellini.top":
		port = "1337"
	case "home.ebinbellini.top":
		port = "4918"
	case "ebin.ebinbellini.top":
		// ころねが踊りだす！
		http.Redirect(w, r, "https://www.youtube.com/watch?v=W9paQ-ZmvbI", http.StatusSeeOther)
		return
	default:
		port = "9001"
	}

	// Create localhost path
	path := r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		path = path + "/?" + r.URL.RawQuery
	}
	resourcePathArray := []string{"http://localhost:", port, path}
	resourcePath := strings.Join(resourcePathArray, "")

	// Construct a copy of the request
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ETT FEL UPPSTOD", err)
	}
	requestBodyReader := bytes.NewReader(requestBody)
	newRequest, err := http.NewRequest(r.Method, resourcePath, requestBodyReader)
	if err != nil {
		fmt.Println("Nu blev de fel", err)
	}

	// Send request
	res, err := netClient.Do(newRequest)
	if err != nil {
		fmt.Println("DE FEL", err)
		return
	}
	defer res.Body.Close()

	// Respond with response from localhost
	content, err := ioutil.ReadAll(res.Body)
	contentReader := bytes.NewReader(content)
	// TODO use status code from local response
	http.ServeContent(w, r, resourcePath, time.Now(), contentReader)
}
