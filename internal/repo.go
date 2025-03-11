package internal

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5"
)

type CommitFilterFunc func(string, TypeValue[ChangelogConditionType, []string]) bool

func OpenRepository(repository string) (*git.Repository, error) {
	r, err := git.PlainOpen(repository)
	if err != nil {
		return nil, fmt.Errorf("repository [%s]: %w", repository, err)
	}

	return r, nil
}

func ChangelogList(repository string, rules Changelog) ([]string, error) {
	if rules.From.Type == "" || rules.To.Type == "" {
		return []string{}, nil
	}

	if rules.From.Value == "" || rules.To.Value == "" {
		return []string{}, nil
	}

	r, err := OpenRepository(repository)
	if err != nil {
		return nil, err
	}

	commits, err := listOfCommits(
		r,
		rules,
		CommitFilter,
	)
	if err != nil {
		return nil, err
	}

	if rules.Sort == Asc {
		slices.Sort(commits)
	}

	if rules.Sort == Desc {
		slices.Sort(commits)
		slices.Reverse(commits)
	}

	return commits, nil
}

func CommitFilter(message string, conditions TypeValue[ChangelogConditionType, []string]) bool {
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
			if conditions.Type == Include {
				return true
			}

			if conditions.Type == Exclude {
				return false
			}
		}

		matched = !(conditions.Type == Include)
	}

	return matched
}

func listOfCommits(
	repository *git.Repository,
	rules Changelog,
	filter CommitFilterFunc,
) ([]string, error) {
	if repository == nil {
		return nil, errors.New("repository is nil")
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

	result := make([]string, 0)
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

func hashes(repository *git.Repository, rules Changelog) (plumbing.Hash, plumbing.Hash, error) {
	if repository == nil {
		return plumbing.ZeroHash, plumbing.ZeroHash, errors.New("repository is nil")
	}

	var startHash plumbing.Hash
	var endHash plumbing.Hash

	if rules.From.Type == Commit {
		startHash = plumbing.NewHash(rules.From.Value)
	} else {
		hash, err := repository.ResolveRevision(plumbing.Revision(rules.From.Value))
		if err != nil {
			return startHash, endHash, fmt.Errorf("failed to resolve commit hash [%s]: %w", rules.From.Value, err)
		}
		startHash = *hash
	}

	if rules.To.Type == Commit {
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
