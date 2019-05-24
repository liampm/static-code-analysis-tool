package domain

import (
	uuid "github.com/satori/go.uuid"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type TaskInstance struct {
	Id        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"projectId"`
	Name      string    `json:"name"`
}

type task interface {
	analyse(target *Target) interface{}
}

type fileCount struct {
}

func (task *fileCount) analyse(target *Target) interface{} {
	count := 0

	directory := target.directory("fc")

	_, err := cloneRepo(target.Config.Details.(RepoDetails), directory)
	handleErr(directory, err)

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {


		if !info.IsDir() {
			count++
		}
		return nil
	})

	handleErr(directory, err)

	err = os.RemoveAll(directory)

	if err != nil {
		log.Printf("Failed to clear up output directory '%s'", directory)
	}

	type Result struct {
		Count int `json:"count"`
	}

	return Result{
		Count: count,
	}
}

type addAFile struct {
	filename string
	content string
}

func (task *addAFile) analyse(target *Target) interface{} {
	directory := target.directory("ae")

	repo, err := cloneRepo(target.Config.Details.(RepoDetails), directory)
	handleErr(directory, err)

	worktree, err := repo.Worktree()
	handleErr(directory, err)


	headRef, err := repo.Head()
	handleErr(directory, err)

	branchName := plumbing.NewBranchReferenceName("my-new-branch" + uuid.NewV4().String())

	ref := plumbing.NewHashReference(branchName, headRef.Hash())


	err = repo.Storer.SetReference(ref)
	handleErr(directory, err)


	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchName,
		Create: false,
	})
	handleErr(directory, err)


	filename := filepath.Join(directory, task.filename)
	err = ioutil.WriteFile(filename, []byte(task.content), 0644)
	handleErr(directory, err)

	_, err = worktree.Add(task.filename)
	handleErr(directory, err)


	commit, err := worktree.Commit("Change from SCAT", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "SCAT",
			Email: "SCAT@maruhub.com",
			When:  time.Now(),
		},
	})
	handleErr(directory, err)

	_ = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: target.Config.Details.(RepoDetails).Username,
			Password: target.Config.Details.(RepoDetails).Token,
		},
	})

	err = os.RemoveAll(directory)

	if err != nil {
		log.Printf("Failed to clear up output directory '%s'", directory)
	}

	return map[string]string{
		"commit": commit.String(),
	}
}


func cloneRepo(config RepoDetails, dir string) (*git.Repository, error) {
	log.Printf("Loading repo into %s", dir)

	return git.PlainClone(dir, false, &git.CloneOptions{
		URL:config.Url,
		Auth: &http.BasicAuth{
			Username: config.Username,
			Password: config.Token,
		},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
}


func handleErr (directory string ,err error) {
	if err != nil {
		err2 := os.RemoveAll(directory)

		if err2 != nil {
			log.Printf("Failed to clear up output directory '%s'", directory)
		}
		panic(err)
	}
}