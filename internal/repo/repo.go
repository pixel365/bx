package repo

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5"
)

var (
	listOfCommitsFunc  = listOfCommits
	openRepositoryFunc = OpenRepository
	hashesFunc         = hashes
)

type CommitFilterFunc func(string, types.TypeValue[types.ChangelogConditionType, []string]) bool

// OpenRepository opens an existing Git repository at the specified path.
//
// Parameters:
//   - repository: The file system path to the Git repository.
//
// Returns:
//   - A pointer to a `git.Repository` instance if the repository is successfully opened.
//   - An error if the repository cannot be opened.
//
// Behavior:
//   - Uses `git.PlainOpen` to open the repository.
//   - Wraps the error with additional context if opening fails.
//
// Example:
//
//	repo, err := OpenRepository("/path/to/repo")
//	if err != nil {
//	    log.Fatalf("Failed to open repository: %v", err)
//	}
//	fmt.Println("Repository opened:", repo)
//
// Notes:
//   - The function does not initialize a new repository; it only opens an existing one.
//   - Ensure that the provided path is a valid Git repository.
func OpenRepository(repository string) (*git.Repository, error) {
	r, err := git.PlainOpen(repository)
	if err != nil {
		return nil, fmt.Errorf("repository [%s]: %w", repository, err)
	}

	return r, nil
}

// ChangelogList generates a list of commits between two specified points in a Git repository,
// applying given changelog rules.
//
// Parameters:
//   - repository: The file system path to the Git repository.
//   - rules: A `Changelog` struct defining the range of commits and sorting rules.
//
// Returns:
//   - A slice of commit hashes as strings.
//   - An error if the repository cannot be opened or commit retrieval fails.
//
// Behavior:
//   - If `rules.From` or `rules.To` are not properly set, it returns an empty list with no error.
//   - Opens the specified Git repository using `OpenRepository`.
//   - Retrieves commits within the specified range using `listOfCommits`.
//   - Sorts commits if `rules.Sort` is set to `Asc` (ascending) or `Desc` (descending).
//
// Example:
//
//	rules := Changelog{
//	    From: Tag{"v1.0.0"},
//	    To:   Tag{"v2.0.0"},
//	    Sort: Desc,
//	}
//	commits, err := ChangelogList("/path/to/repo", rules)
//	if err != nil {
//	    log.Fatalf("Failed to generate changelog: %v", err)
//	}
//	fmt.Println("Commits:", commits)
//
// Notes:
//   - Sorting is applied only if `rules.Sort` is explicitly set.
//   - Uses `listOfCommits` with a predefined `CommitFilter` function.
//   - The function does not modify the repository; it only queries commit history.
func ChangelogList(repository string, rules types.Changelog) ([]string, error) {
	if rules.From.Type == "" || rules.To.Type == "" {
		return []string{}, nil
	}

	if rules.From.Value == "" || rules.To.Value == "" {
		return []string{}, nil
	}

	r, err := openRepositoryFunc(repository)
	if err != nil {
		return nil, err
	}

	commits, err := listOfCommitsFunc(
		r,
		rules,
		CommitFilter,
	)
	if err != nil {
		return nil, err
	}

	if rules.Sort == types.Asc {
		slices.Sort(commits)
	}

	if rules.Sort == types.Desc {
		slices.Sort(commits)
		slices.Reverse(commits)
	}

	return commits, nil
}

// CommitFilter checks whether a commit message matches a set of conditions.
//
// Parameters:
//   - message: The commit message to evaluate.
//   - conditions: A `TypeValue` struct containing a list of regex conditions and a filter type
//     (`Include` or `Exclude`).
//
// Returns:
//   - true if the commit message satisfies the filtering conditions.
//   - false if the message does not match an `Include` condition or matches an `Exclude` condition.
//
// Behavior:
//   - If `conditions.Value` is empty, it returns `true` (no filtering applied).
//   - Iterates through all regex conditions in `conditions.Value`.
//   - If `conditions.Type` is `Include`, returns `true` on the first match.
//   - If `conditions.Type` is `Exclude`, returns `false` on the first match.
//   - If no matches are found, returns `false` for `Include` and `true` for `Exclude`.
//
// Example:
//
//	conditions := TypeValue[ChangelogConditionType, []string]{
//	    Type:  Include,
//	    Value: []string{".*fix.*", ".*feature.*"},
//	}
//	fmt.Println(CommitFilter("fix: bug fixed", conditions)) // true
//	fmt.Println(CommitFilter("chore: update docs", conditions)) // false
//
// Notes:
//   - Regular expressions are compiled dynamically; invalid regex patterns are ignored.
//   - If regex compilation fails, the loop breaks, and filtering may be incomplete.
//   - Ensures `matched` is set correctly for both `Include` and `Exclude` conditions.
func CommitFilter(
	message string,
	conditions types.TypeValue[types.ChangelogConditionType, []string],
) bool {
	if len(conditions.Value) == 0 {
		return true
	}

	matched := true
	for _, condition := range conditions.Value {
		reg, err := regexp.Compile(condition)
		if err != nil {
			break
		}

		if reg.MatchString(message) {
			if conditions.Type == types.Include {
				return true
			}

			if conditions.Type == types.Exclude {
				return false
			}
		}

		matched = !(conditions.Type == types.Include)
	}

	return matched
}

// listOfCommits retrieves a list of commit messages from a Git repository
// within a specified range, applying a filtering function.
//
// Parameters:
//   - repository: A pointer to a `git.Repository` instance.
//   - rules: A `Changelog` struct defining the commit range and filtering rules.
//   - filter: A function that determines whether a commit message should be included.
//
// Returns:
//   - A slice of commit messages that match the filtering criteria.
//   - An error if the repository is nil, commit history retrieval fails, or iteration encounters an issue.
//
// Behavior:
//   - Calls `hashes` to determine the start and end commit hashes based on `rules`.
//   - Retrieves the commit log starting from `endHash`.
//   - Iterates through commits, applying `filter` to the first line of each commit message.
//   - Stops iteration when the `startHash` commit is reached.
//   - If a commit matches the filter condition, its message is added to the result slice.
//
// Example:
//
//	commits, err := listOfCommits(repo, rules, CommitFilter)
//	if err != nil {
//	    log.Fatalf("Failed to list commits: %v", err)
//	}
//	fmt.Println("Filtered commits:", commits)
//
// Notes:
//   - Uses `plumbing.ErrObjectNotFound` to stop processing when `startHash` is reached.
//   - Ensures the commit iterator is closed using `defer iter.Close()`.
//   - If an error occurs while iterating, it is wrapped and returned unless it's `ErrObjectNotFound`.
func listOfCommits(
	repository *git.Repository,
	rules types.Changelog,
	filter CommitFilterFunc,
) ([]string, error) {
	if repository == nil {
		return nil, errors2.NilRepositoryError
	}

	startHash, endHash, err := hashes(repository, rules)
	if err != nil {
		return nil, fmt.Errorf("repository [%s]: %w", repository, err)
	}

	iter, err := repository.Log(&git.LogOptions{From: endHash})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve commit history: %w", err)
	}
	defer iter.Close()

	var result []string
	err = iter.ForEach(func(c *object.Commit) error {
		if c.Hash == startHash {
			return plumbing.ErrObjectNotFound
		}

		m := strings.Split(c.Message, "\n")
		if filter(m[0], rules.Condition) {
			result = append(result, m[0])
		}
		return nil
	})

	if err != nil && !errors.Is(err, plumbing.ErrObjectNotFound) {
		return nil, fmt.Errorf("failed to iterate commit history: %w", err)
	}

	return result, nil
}

// hashes resolves the start and end commit hashes based on the given changelog rules.
//
// Parameters:
//   - repository: A pointer to a `git.Repository` instance.
//   - rules: A `Changelog` struct specifying the commit range.
//
// Returns:
//   - `startHash`: The resolved plumbing.Hash for the starting commit.
//   - `endHash`: The resolved plumbing.Hash for the ending commit.
//   - An error if the repository is nil or if a commit hash cannot be resolved.
//
// Behavior:
//   - If the repository is nil, returns `plumbing.ZeroHash` for both values and an error.
//   - If `rules.From.Type` is `Commit`, directly converts `rules.From.Value` into a hash.
//   - Otherwise, attempts to resolve `rules.From.Value` as a branch, tag, or other reference.
//   - Performs the same logic for `rules.To`.
//   - If resolving a commit hash fails, an error is returned.
//
// Example:
//
//	start, end, err := hashes(repo, rules)
//	if err != nil {
//	    log.Fatalf("Failed to resolve commit hashes: %v", err)
//	}
//	fmt.Println("Start Hash:", start, "End Hash:", end)
//
// Notes:
//   - The function supports resolving both direct commit hashes and references (branches/tags).
//   - Uses `repository.ResolveRevision` to translate references into commit hashes.
func hashes(
	repository *git.Repository,
	rules types.Changelog,
) (plumbing.Hash, plumbing.Hash, error) {
	if repository == nil {
		return plumbing.ZeroHash, plumbing.ZeroHash, errors2.NilRepositoryError
	}

	var startHash plumbing.Hash
	var endHash plumbing.Hash

	if rules.From.Type == types.Commit {
		startHash = plumbing.NewHash(rules.From.Value)
	} else {
		hash, err := repository.ResolveRevision(plumbing.Revision(rules.From.Value))
		if err != nil {
			return startHash, endHash, fmt.Errorf("failed to resolve commit hash [%s]: %w", rules.From.Value, err)
		}
		startHash = *hash
	}

	if rules.To.Type == types.Commit {
		endHash = plumbing.NewHash(rules.To.Value)
	} else {
		hash, err := repository.ResolveRevision(plumbing.Revision(rules.To.Value))
		if err != nil {
			return startHash, endHash, fmt.Errorf("failed to resolve commit hash [%s]: %w", rules.To.Value, err)
		}
		endHash = *hash
	}

	return startHash, endHash, nil
}

// ChangesList generates a list of file changes between two commits in a Git repository.
//
// Parameters:
//   - repository: The file system path to the Git repository.
//   - rules: A `Changelog` struct defining the commit range for comparison.
//
// Returns:
//   - A pointer to a `Changes` struct containing lists of added, modified, and deleted files.
//   - An error if the repository cannot be opened, commit hashes cannot be resolved, or patch generation fails.
//
// Behavior:
//   - Opens the Git repository using `OpenRepository`.
//   - Resolves the start and end commit hashes using `hashes`.
//   - Retrieves the commit objects corresponding to these hashes.
//   - Generates a diff (`Patch`) between the two commits.
//   - Iterates through the `FilePatches()` to categorize files as added, modified, or deleted.
//
// Example:
//
//	changes, err := ChangesList("/path/to/repo", rules)
//	if err != nil {
//	    log.Fatalf("Failed to generate changes list: %v", err)
//	}
//	fmt.Println("Added files:", changes.Added)
//	fmt.Println("Modified files:", changes.Modified)
//	fmt.Println("Deleted files:", changes.Deleted)
//
// Notes:
//   - If `from == nil`, the file was newly added.
//   - If `to == nil`, the file was deleted.
//   - Otherwise, the file was modified.
//   - The function does not modify the repository; it only analyzes commit differences.
func ChangesList(repository string, rules types.Changelog) (*types.Changes, error) {
	r, err := openRepositoryFunc(repository)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, errors2.NilRepositoryError
	}

	startHash, endHash, err := hashesFunc(r, rules)
	if err != nil {
		return nil, fmt.Errorf("repository [%s]: %w", repository, err)
	}

	startCommit, err := r.CommitObject(startHash)
	if err != nil {
		return nil, fmt.Errorf("repository [%s]: %w", repository, err)
	}

	endCommit, err := r.CommitObject(endHash)
	if err != nil {
		return nil, fmt.Errorf("repository [%s]: %w", repository, err)
	}

	patch, err := startCommit.Patch(endCommit)
	if err != nil {
		return nil, fmt.Errorf("repository [%s]: %w", repository, err)
	}

	c := types.Changes{}

	for _, filePatch := range patch.FilePatches() {
		from, to := filePatch.Files()

		if from == nil && to != nil {
			c.Added = append(c.Added, to.Path())
		}

		if from != nil && to == nil {
			c.Deleted = append(c.Deleted, from.Path())
		}

		if from != nil && to != nil {
			if from.Path() != to.Path() {
				c.Moved = append(c.Moved, to.Path())
			} else {
				c.Modified = append(c.Modified, from.Path())
			}
		}
	}

	return &c, nil
}
