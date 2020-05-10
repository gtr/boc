package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func makeAuthRequest(requestType, url string) []byte {
	// Request the Github API.
	request, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Set the Github personal access token.
	token := getAccessToken()
	request.SetBasicAuth(token, "x-oauth-basic")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("StatusCode:", response.StatusCode)

	// Read the body of the request.
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(body))

	return body
}

func getAccessToken() string {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var result map[string]string
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)

	return result["MY_KEY"]
}

func clone(url string) {
	cmd := exec.Command("git", "clone", url)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func remove(dir string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(pwd)
	files, err := filepath.Glob(filepath.Join(pwd, dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}
