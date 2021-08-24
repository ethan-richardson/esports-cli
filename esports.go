package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type matchInfo struct {
	Matchup    string    `json:"name"`
	Start      time.Time `json:"begin_at"`
	Tournament struct {
		Name string `json:"name"`
	} `json:"tournament"`
	League struct {
		Name string `json:"name"`
	} `json:"league"`
}

func getData(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var match []matchInfo

	err = json.Unmarshal(body, &match)
	if err != nil {
		log.Fatal(err)
	}
	processJson(match)
}

func createRequest(game string, team string, daysAhead int) string {
	teamAPIString := "&search[name]=" + team
	file, err := os.Open("token.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	token := "&token=" + scanner.Text()

	if daysAhead > 0 {
		currentTime := time.Now().UTC().Format(time.RFC3339)
		endTime := time.Now().AddDate(0, 0, daysAhead).UTC().Format(time.RFC3339)
		timeAPIString := "&range[begin_at]=" + currentTime + "," + endTime
		return "https://api.pandascore.co/" + game + "/matches/upcoming?" + teamAPIString + timeAPIString + token
	} else {
		return "https://api.pandascore.co/" + game + "/matches/upcoming?" + teamAPIString + token
	}
}

func processJson(matchList []matchInfo) {
	fmt.Println("Matches:")
	for i := range matchList {
		fmt.Println(matchList[i].League.Name)
		// fmt.Println(matchList[i].Tournament.Name)
		fmt.Println(matchList[i].Matchup)
		fmt.Println(matchList[i].Start.Local().Format("Mon Jan _2 03:00pm"), "\n")
	}
}

func main() {
	game := os.Args[1]
	team := os.Args[2]
	var url string
	if len(os.Args) > 3 {
		daysAhead, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		url = createRequest(game, team, daysAhead)
	} else {
		url = createRequest(game, team, 0)
	}
	getData(url)
}
