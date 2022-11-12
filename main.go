package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var fruits map[string]int

func main() {
	fruits = make(map[string]int)
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			respond(w)
		})
	http.HandleFunc("/buy", buy)
	http.HandleFunc("/sell", sell)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "9999"
	}
	fmt.Println("listening on port: ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func buy(w http.ResponseWriter, r *http.Request) {
	fruit := r.URL.Query().Get("fruit")
	if len(fruit) == 0 {
		respondWithError(w, "parameter fruit is required\n")
		return
	}
	count := r.URL.Query().Get("count")
	if len(count) == 0 {
		respondWithError(w, "parameter count is required\n")
		return
	}
	c, err := strconv.Atoi(count)
	if err != nil {
		respondWithError(w, "count must be a number\n")
		return
	}
	fruits[fruit] += c
	respond(w)
}

func sell(w http.ResponseWriter, r *http.Request) {
	fruit := r.URL.Query().Get("fruit")
	if len(fruit) == 0 {
		respondWithError(w, "parameter fruit is required\n")
		return
	}
	count := r.URL.Query().Get("count")
	if len(count) == 0 {
		respondWithError(w, "count is required\n")
		return
	}
	c, err := strconv.Atoi(count)
	if err != nil {
		respondWithError(w, "count must be a number\n")
		return
	}
	if c > fruits[fruit] {
		respondWithError(w, "not enough fruits to sell\n")
		return
	}
	fruits[fruit] -= c
	respond(w)
}

func respond(w http.ResponseWriter) {
	b, _ := yaml.Marshal(fruits)
	if strings.TrimSpace(string(b)) == "{}" {
		fmt.Fprint(w, "no fruits\n")
		return
	}
	fmt.Fprint(w, string(b))
}

func respondWithError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(message))
}
