package main

import (
	"fmt"
	"github.com/BasixKOR/telescope-updater/utils"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading .env file")
		os.Exit(1)
	}

	fmt.Println("Initializing Agolia...")
	agolia := search.NewClient(GetKey("AGOLIA_APP"), GetKey("AGOLIA_KEY"))
	index := agolia.InitIndex("repos")

	fmt.Println("Initializing GitHub API...")
	key := GetKey("GITHUB_TOKEN")
	github := githubv4.NewClient(utils.NewBearerClient(key))

	fmt.Println("Initialized! Attempting to fetch...")
	c := make(chan utils.FetchedRepos)
	go utils.Fetch(github, c)
	for repos := range c {
		_, err := index.SaveObjects(repos)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error saving objects: ", err)
			os.Exit(1)
		}
	}

	fmt.Println("Everything done!")
}

func GetKey(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Print("Put your %s: ", strings.Title(strings.ReplaceAll(key, "_", " ")))
		fmt.Scanln(&value)
	}
	value = strings.Trim(value, " ")
	return value
}
