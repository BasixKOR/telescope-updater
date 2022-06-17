package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
)

type StarredReposQuery struct {
	User struct {
		StarredRepositories struct {
			Nodes []struct {
				ID            githubv4.ID
				NameWithOwner string
				Description   string
				Stargazers    struct {
					TotalCount int
				}
				PrimaryLanguage struct {
					Name string
				}
				RepositoryTopics struct {
					Nodes []struct {
						Topic struct {
							Name string
						}
					}
				} `graphql:"repositoryTopics(first: 30)"`
			}
			PageInfo PageInfo
		} `graphql:"starredRepositories(first: 100, after: $page)"`
	} `graphql:"user(login: \"BasixKOR\")"`
}

type PageInfo struct {
	EndCursor   string
	HasNextPage bool
}

type FetchedRepo struct {
	ID            githubv4.ID `json:"objectID,string"`
	NameWithOwner string      `json:"nameWithOwner"`
	Description   string      `json:"description"`
	Stargazers    struct {
		TotalCount int `json:"totalCount"`
	} `json:"stargazers"`
	PrimaryLanguage struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	Topics []string `json:"topics"`
}

func convert(queried StarredReposQuery) []FetchedRepo {
	original := queried.User.StarredRepositories.Nodes
	fetched := []FetchedRepo{}
	for _, i := range original {
		topics := []string{}

		for _, n := range i.RepositoryTopics.Nodes {
			topics = append(topics, n.Topic.Name)
		}

		fetched = append(fetched, FetchedRepo{
			ID:            i.ID,
			NameWithOwner: i.NameWithOwner,
			Description:   i.Description[:200],
			Stargazers: struct {
				TotalCount int `json:"totalCount"`
			}{
				i.Stargazers.TotalCount,
			},
			PrimaryLanguage: struct {
				Name string `json:"name"`
			}{
				i.PrimaryLanguage.Name,
			},
			Topics: topics,
		})

	}
	return fetched
}

var query StarredReposQuery

func Fetch(client *githubv4.Client, c chan []FetchedRepo) {
	pageInfo := map[string]interface{}{
		"page": (*githubv4.String)(nil),
	}
	for i := 1; ; i++ {
		fmt.Println("Requesting chunk:", i)
		err := client.Query(context.Background(), &query, pageInfo)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			close(c)
			return
		}
		c <- convert(query)
		if !query.User.StarredRepositories.PageInfo.HasNextPage {
			break
		}
		pageInfo["page"] = githubv4.String(query.User.StarredRepositories.PageInfo.EndCursor)
	}
	close(c)
}
