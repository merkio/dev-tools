package utils

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CheckoutTo change branch to the @branch
func CheckoutTo(branch string, directory string) {
	fmt.Println("Checkout to branch")

	r, _ := git.PlainOpen(directory)

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		panic(fmt.Errorf("Fatal error get working tree for the repository: \n%s", err))
	}

	// ... checking out to commit
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})

	if err != nil {
		panic(fmt.Errorf("Fatal error checkout to the branch: %s\n%s", branch, err))
	}
}

// CloneRepo clone remote repo to the workspace
func CloneRepo(remoteURL string, directory string) {
	fmt.Println("Clone remote repo")

	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      remoteURL,
		Progress: os.Stdout,
	})

	if err != nil {
		panic(fmt.Errorf("Fatal error during clone repo %s to the directory %s\n%s", remoteURL, directory, err))
	}
}

// IsRepoClean check if repository contains uncommitted changes
func IsRepoClean(directory string) bool {
	r, _ := git.PlainOpen(directory)

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		panic(fmt.Errorf("Fatal error get working tree for the repository: \n%s", err))
	}

	s, err := w.Status()

	if err != nil {
		panic(fmt.Errorf("Fatal error get status for the repository: \n%s", err))
	}
	return s.IsClean()
}

// PullChanges pull the latest changes in the current branch
func PullChanges(repository string) {
	fmt.Println("Pull changes")

	// Open the git repository
	r, _ := git.PlainOpen(repository)

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		panic(fmt.Errorf("Fatal error get working tree for the repository: \n%s", err))
	}

	s, err := w.Status()

	if !s.IsClean() {
		fmt.Printf("Repository contains uncommitted changes. Repository: %s", repository)
		return
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})

	if err != nil {
		fmt.Printf("Fatal error pull changes in [%s]: \n%s", repository, err)
	}
}
