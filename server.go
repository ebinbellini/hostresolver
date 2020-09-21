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
	fmt.Println("HOSTEN ÄR", r.Host)
	// Choose witch port to fetch from
	// TODO tidy up maybe with a switch or map
	var port string
	if r.Host == "ebinbellini.top" || r.Host == "www.ebinbellini.top" {
		port = "9001"
	} else if r.Host == "chat.ebinbellini.top" {
		port = "1337"
	} else if r.Host == "home.ebinbellini.top" {
		port = "4918"
	} else {
		// No match
		// Respond with ころね instead
		http.Redirect(w, r, "https://www.youtube.com/watch?v=W9paQ-ZmvbI", http.StatusSeeOther)
		return
	}

	// Create localhost path
	path := r.URL.Path
	resourcePathArray := []string{"http://localhost:", port, path}
	resourcePath := strings.Join(resourcePathArray, "")

	// Construct new request
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
		fmt.Fprint(w, err)
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
