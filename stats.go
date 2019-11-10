package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"strings"
	"time"
)

func stats() map[int]int {
	const name = "max"

	fp := getDotFilePath()
	repos := parseFileLinesToSlice(fp)

	commits := make(map[int]int, 183)
	for i := 183; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(name, path, commits)
	}

	log.Println(commits)
	return commits
}

func fillCommits(name string, path string, commits map[int]int) map[int]int {
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
		daysAgo := daysSince(c.Author.When)
		if strings.Contains(c.Author.Email, name) && daysAgo != 1000 {
			commits[daysAgo]++
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
