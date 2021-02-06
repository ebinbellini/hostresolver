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
	http.HandleFunc("/", serve)

	fmt.Println("Listening on :443...")
	err := http.ListenAndServeTLS(":443",
		"/etc/letsencrypt/live/ebinbellini.top/fullchain.pem",
		"/etc/letsencrypt/live/ebinbellini.top/privkey.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	ct := time.Now()

	space := " "
	for i := 0; i < 20-len(r.Host); i++ {
		space = space + " "
	}
	fmt.Println(ct.Format("2006-01-02 15:04:05"), "-=- HOST =", r.Host, space, "IP =", r.RemoteAddr)

	port := resolveHostPort(w, r)
	if port == "" {
		return
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
	// TODO use header from local response
	http.ServeContent(w, r, resourcePath, time.Now(), contentReader)
}

func resolveHostPort(w http.ResponseWriter, r *http.Request) string {
	// Choose witch port to fetch from
	switch r.Host {
	case "ebinbellini.top":
		return "9001"
	case "www.ebinbellini.top":
		return "9001"
	case "chat.ebinbellini.top":
		return "1337"
	case "home.ebinbellini.top":
		return "4918"
	case "weather.ebinbellini.top":
		return "737"
	case "ebin.ebinbellini.top":
		// ころねが踊りだす！
		http.Redirect(w, r, "https://www.youtube.com/watch?v=iFlBEnW90oE", http.StatusSeeOther)
		return ""
	default:
		return "9001"
	}
}
