package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var netClient = &http.Client{
	Timeout: 60 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect attempt")
	},
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

	// Redirect from *.top to *.com
	if strings.HasSuffix(r.Host, ".top") {
		redirectToCom(w, r)
		return
	}

	port := resolveHostPort(w, r)
	if port == "" {
		// Request handled by
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
		serveError(w, r)
		return
	}

	requestBodyReader := bytes.NewReader(requestBody)
	newRequest, err := http.NewRequest(r.Method, resourcePath, requestBodyReader)
	if err != nil {
		serveError(w, r)
		return
	}

	// Send request
	// TODO check if error is serious
	res, _ := netClient.Do(newRequest)
	defer res.Body.Close()

	// Respond with response from localhost
	content, err := ioutil.ReadAll(res.Body)
	contentReader := bytes.NewReader(content)

	// Copy all headers
	for name, values := range res.Header {
		w.Header()[name] = values
	}

	// Copy the status code
	w.WriteHeader(res.StatusCode)

	http.ServeContent(w, r, resourcePath, time.Now(), contentReader)
}

func resolveHostPort(w http.ResponseWriter, r *http.Request) string {
	// Choose witch port to fetch from
	switch r.Host {
	case "ebinbellini.com":
		return "9001"
	case "www.ebinbellini.com":
		return "9001"
	case "chat.ebinbellini.com":
		return "1337"
	case "home.ebinbellini.com":
		return "4918"
	case "weather.ebinbellini.com":
		return "737"
	case "ebin.ebinbellini.com":
		// ころねが踊りだす！
		http.Redirect(w, r, "https://www.youtube.com/watch?v=iFlBEnW90oE", http.StatusSeeOther)
		return ""
	default:
		return "9001"
	}
}

func serveError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	http.ServeFile(w, r, "error.html")
}

func redirectToCom(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`\.top`)
	newHost := re.ReplaceAllString(r.Host, ".com")
	http.Redirect(w, r, "https://"+newHost+r.URL.Path, http.StatusMovedPermanently)
}
