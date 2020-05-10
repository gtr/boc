package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// API base
const (
	api   = "https://api.github.com"
	users = "/users"
	repos = "/repos"
	lang  = "/languages"
	get   = "GET"
)

// langPair represents a pair of language and its associated byte count.
type langPair struct {
	Lang  string
	Count int
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		return
	}
	repoHTTPBody := getAllRepositories(os.Args[1])
	allRepositories := parseForRepositories(repoHTTPBody)
	// printRepositories(allRepositories)
	langCount := countLanguages(allRepositories)
	sortLanguages(langCount)

}

// getAllRepositories makes an authenticated HTTP request to find all the
// public repositories for a given user.
func getAllRepositories(user string) []byte {
	// Request the Github API.
	url := api + users + "/" + user + repos
	return makeAuthRequest(get, url)
}

// printRepositories prints all the repositories line-by-line.
func printRepositories(allRepositories [1000]string) {
	for _, repo := range allRepositories {
		if len(repo) != 0 {
			fmt.Println(repo)
		}
	}
}

// parseForRepositories takes in an response body to parse for all repositories.
func parseForRepositories(body []byte) [1000]string {
	search := regexp.MustCompile("\"full_name\":\"")
	bodyString := string(body)
	matches := search.FindAllStringIndex(bodyString, -1)

	var allRepositories [1000]string

	for i, pair := range matches {
		start := pair[1]
		search := bodyString[start:]
		end := strings.Index(search, "\"")
		currRepository := bodyString[start : start+end]
		allRepositories[i] = currRepository

	}

	return allRepositories
}

// countLanguages queries the Githb API to find the number of bytes for each
// language.
func countLanguages(allRepositories [1000]string) map[string]int {
	// Initialize language count map.
	var langCount map[string]int
	langCount = make(map[string]int)

	for _, repo := range allRepositories {
		if len(repo) != 0 {
			var result map[string]int
			// Request the Github APi.
			url := api + repos + "/" + repo + lang
			body := makeAuthRequest(get, url)

			json.Unmarshal([]byte(body), &result)

			for lang, count := range result {
				if val, ok := langCount[lang]; ok {
					langCount[lang] = count + val
				} else {
					langCount[lang] = count
				}
			}
		}
	}

	// fmt.Println(langCount)
	return langCount
}

// sortLanguages sorts the languages based on the number of bytes.
func sortLanguages(langCount map[string]int) []langPair {
	var sortedSlice []langPair
	for lang, count := range langCount {
		sortedSlice = append(sortedSlice, langPair{lang, count})
	}

	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].Count > sortedSlice[j].Count
	})

	for _, pair := range sortedSlice {
		fmt.Printf("%s -> %d\n", pair.Lang, pair.Count)
	}

	return sortedSlice
}

func printUsage() {
	fmt.Printf("Error: incorrect usage\n\nUsage: ./boc [USER]\n\nArguments:\n")
	fmt.Printf("  USER\tThe username of the github user to query.\n")
}
