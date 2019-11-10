package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

// check if we need to scan for new repositories
func checkForPathList() bool {
	if _, err := os.Stat(getDotFilePath()); os.IsNotExist(err) {
		return false
	}
	return true
}

// scan scans a new folder for Git repositories
func scan() {
	repositories := recursiveScanFolder(getHomeDirPath())
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repositories)
}

// recursiveScanFolder starts the recursive search of git repositories
// living in the `folder` subtree
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

// scanGitFolders returns a list of subfolders of `folder` ending with `.git`.
// Returns the base folder of the repo, the .git folder parent.
// Recursively searches in the subfolders by passing an existing `folders` slice.
func scanGitFolders(folders []string, folder string) []string {
	// trim the last `/`
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal("line 43 ", err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal("line 48 ", err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" || file.Name() == ".oh-my-zsh" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}

	return folders
}

func getHomeDirPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("line 75 ", err)
	}
	return usr.HomeDir
}

// getDotFilePath returns the dot file for the repos list
func getDotFilePath() string {
	return getHomeDirPath() + "/.commitment"
}

// addNewSliceElementsToFile given a slice of strings representing paths, stores them
// to the filesystem
func addNewSliceElementsToFile(filePath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

// parseFileLinesToSlice given a file path string, gets the content
// of each line and parses it to a slice of strings.
func parseFileLinesToSlice(filePath string) []string {
	var lines []string

	file, err := openFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("line 120 ", err)
	}

	return lines
}

// openFile opens the file located at `filePath`. Creates it if not existing.
func openFile(filePath string) (*os.File, error) {
	var file *os.File

	if _, err := os.Stat(getDotFilePath()); os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			return file, err
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, os.ModeAppend)
		if err != nil {
			return file, err
		}
	}

	return file, nil
}

// joinSlices adds the element of the `new` slice
// into the `existing` slice, only if not already there
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

// sliceContains returns true if `slice` contains `value`
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// dumpStringsSliceToFile writes content to the file in path `filePath` (overwriting existing content)
func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	ioutil.WriteFile(filePath, []byte(content), 0755)
}
