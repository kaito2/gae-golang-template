package main

import (
	"encoding/json"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/kaito2/gae-golang-template/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", indexHandler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	// [END setting_port]
}

// ref.
// https://github.com/GoogleCloudPlatform/golang-samples/tree/master/appengine/go11x/helloworld
// https://qiita.com/mztnnrt/items/7a8092e30a54234a912a
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// validate path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// validate HTTP request method
	if r.Method != http.MethodGet && r.Method != http.MethodPost{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "GET and POST method is allowed")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// ref. https://golangcode.com/get-a-url-parameter-from-a-request/
		queryStrings := r.URL.Query()
		name, ok := queryStrings["name"]
		if !ok || len(name) < 1 {
			name = []string{"UNKNOWN"}
		}
		age, ok := queryStrings["age"]
		if !ok || len(age) < 1 {
			age = []string{"0"}
		}
		fmt.Fprintf(w,"(POST) Your name: %s, Your Age: %s", name, age)
	case http.MethodPost:
		length, err := strconv.Atoi(r.Header.Get("Content-Length"))
		if err != nil {
			log.Print("Content-LengthがnilだからStatusInternalServerError")
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Print("lengthは", length, "バイトです。")
		}

		var user models.User
		buffer := make([]byte, length)

		_, err = r.Body.Read(buffer)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(buffer, &user)
		// use third party package
		myFigure := figure.NewFigure(fmt.Sprintf("Hello %s", user.Name), "", true)
		fmt.Fprintf(w, "%s\n", myFigure)
		fmt.Fprintf(w,"(POST) Your name: %s, Your Age: %d\n", user.Name, user.Age)
	}
}
