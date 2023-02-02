package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var fruits map[string]int

type Req struct {
	Fruit string `json:"fruit"`
	Count int    `json:"count"`
}

type Res struct {
	Fruits map[string]int `json:"fruits"`
}

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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
}

func getReq(w http.ResponseWriter, r *http.Request) *Req {
	var req Req
	fruit := r.URL.Query().Get("fruit")
	countStr := r.URL.Query().Get("count")
	fromParams := true
	if len(fruit) == 0 || len(countStr) == 0 {
		fromParams = false
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			respondWithError(w, "either fruit or count is not provided but required")
			return nil
		}
	}
	if fromParams {
		count, err := strconv.Atoi(countStr)
		if err != nil {
			respondWithError(w, "count must be a number\n")
			return nil
		}
		req.Fruit = fruit
		req.Count = count
	}
	if len(req.Fruit) == 0 || req.Count == 0 {
		respondWithError(w, "either fruit or count is not provided but required")
		return nil
	}
	if req.Count <= 0 {
		respondWithError(w, "count must be a positive number\n")
		return nil
	}
	return &req
}

func buy(w http.ResponseWriter, r *http.Request) {
	req := getReq(w, r)
	if req == nil {
		return
	}
	fruit := req.Fruit
	c := req.Count
	fruits[fruit] += c
	respond(w)
}

func sell(w http.ResponseWriter, r *http.Request) {
	req := getReq(w, r)
	if req == nil {
		return
	}
	fruit := req.Fruit
	c := req.Count
	if fruits[fruit] < c {
		respondWithError(w, "not enough fruits\n")
		return
	}
	fruits[fruit] -= c
	respond(w)
}

func respond(w http.ResponseWriter) {
	enableCors(&w)
	b, _ := json.Marshal(Res{fruits})
	fmt.Fprint(w, string(b))
}

func respondWithError(w http.ResponseWriter, message string) {
	enableCors(&w)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(message))
}
