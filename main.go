// time-stamp microservice receives timestamps and returns json form of natural/unix timestamp.
// Exercise by Free Code Camp.
// I have to give lots of credit for an invaluable code review!
// Code reviewed by: @groob
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Date struct returned as json.
type Date struct {
	Unix    *int64  `json:"unix"`
	Natural *string `json:"natural"`
}

func main() {
	http.HandleFunc("/", timestamp)

	server := http.Server{
		Addr: ":" + os.Getenv("PORT"),
	}
	fmt.Println("Listening on [", server.Addr, "]...")
	log.Fatal(server.ListenAndServe())
}

func timestamp(w http.ResponseWriter, r *http.Request) {

	tstr := r.URL.Path[1:]

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "utf-8")

	t, err := parseTime(tstr)
	if err != nil {
		// Ignore this error and return nil json.
		if err = json.NewEncoder(w).Encode(Date{nil, nil}); err != nil {
			http.Error(w, "Json string failed to marshal", http.StatusInternalServerError)
		}
		return
	}

	naturalDate := t.Format("January 2, 2006")
	naturalTime := t.Unix()

	date := Date{
		Unix:    &naturalTime,
		Natural: &naturalDate}

	if err := json.NewEncoder(w).Encode(&date); err != nil {
		http.Error(w, "Json string failed to marshal", http.StatusInternalServerError)
		return
	}

}

// returnTime takes a string and returns the time.
// Can handle UNIX seconds input and general time in the form:
// "December 15, 2015"
func parseTime(t string) (time.Time, error) {

	// Optimistically try to parse Unix string
	if unixTime, err := strconv.Atoi(t); err == nil {
		return (time.Unix(int64(unixTime), 0)).UTC(), nil
	}

	// Assume general date form.
	genTime := strings.Split(t, " ")
	if len(genTime) != 3 || len(genTime[0]) < 3 {
		return time.Time{}, fmt.Errorf("general form of date not broken into 3 elements: %v", genTime)
	}

	// Re-arrange general time []string to allow time package to parse it.
	timeString := fmt.Sprintf("%s-%s-%02s", genTime[2], genTime[0][:3], genTime[1][:len(genTime[1])-1])

	parsedTime, err := time.Parse("2006-Jan-02", timeString)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse general time: %s", err)
	}

	return parsedTime.UTC(), nil
}
