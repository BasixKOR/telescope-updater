package utils

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"os"
)

type StarredReposQuery struct {
	User struct {
		StarredRepositories struct {
			Nodes    FetchedRepos
			PageInfo PageInfo
		} `graphql:"starredRepositories(first: 100, after: $page)"`
	} `graphql:"user(login: \"BasixKOR\")"`
}

type PageInfo struct {
	EndCursor   string
	HasNextPage bool
}

type FetchedRepos []struct {
	ID            githubv4.ID `json:"objectID,string"`
	NameWithOwner string      `json:"nameWithOwner"`
	Description   string      `json:"description"`
	Stargazers    struct {
		TotalCount int `json:"totalCount"`
	} `json:"stargazers"`
	PrimaryLanguage struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
}

var query StarredReposQuery

func Fetch(client *githubv4.Client, c chan FetchedRepos) {
	pageInfo := map[string]interface{}{
		"page": (*githubv4.String)(nil),
	}
	for i := 1; ; i++ {
		fmt.Println("Requesting chunk:", i)
		err := client.Query(context.Background(), &query, pageInfo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			close(c)
			return
		}
		c <- query.User.StarredRepositories.Nodes
		if !query.User.StarredRepositories.PageInfo.HasNextPage {
			break
		}
		pageInfo["page"] = githubv4.String(query.User.StarredRepositories.PageInfo.EndCursor)
	}
	close(c)
}
