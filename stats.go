package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"strings"
	"time"
)

type Day struct{
	Date time.Time `json:"date"`
	Count int `json:"count"`
}

func stats() []Day {
	name := "max"
	fp := getDotFilePath()
	repos := parseFileLinesToSlice(fp)

	commits := initSlice()

	for _, path := range repos {
		commits = fillCommits(name, path, commits)
	}

	return commits
}

func fillCommits(name string, path string, commits []Day) []Day {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal("line 31 ", err)
	}

	ref, err := repo.Head()
	if err != nil {
		log.Fatal("line 36", err)
	}

	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal("line 41 ", err)
	}

	err = iterator.ForEach(func(c *object.Commit) error {
		if strings.Contains(c.Author.Email, name) && daysSince(c.Author.When) != 1000 {
			i := mapDateToIndex(c.Author.When)
			commits[i].Count++
		}
		return nil
	})
	if err != nil {
		log.Println("path: ", path)
		log.Fatal("line 52 ", err)
	}

	return commits
}

func beginningOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func daysSince(d time.Time) int {
	days := 0
	now := beginningOfDay(time.Now())
	for d.Before(now) {
		d = d.Add(time.Hour * 24)
		days++
		if days > 183 {
			return 1000
		}
	}
	return days
}

func initSlice() []Day {
	slc := make([]Day, 183)
	slc[0].Date = time.Now().AddDate(0, 0, -182)
	for i := 1; i < 183; i++ {
		slc[i].Date = slc[i-1].Date.AddDate(0, 0, 1)
	}
	return slc
}

func mapDateToIndex(date time.Time) int {
	return 182 - daysSince(date)
}

