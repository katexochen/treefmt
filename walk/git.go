package walk

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5/plumbing/format/index"

	"github.com/go-git/go-git/v5"
)

type gitWalker struct {
	root  string
	paths chan string
	repo  *git.Repository
}

func (g *gitWalker) Root() string {
	return g.root
}

func (g *gitWalker) Walk(ctx context.Context, fn WalkFunc) error {
	// for quick relative paths
	relPathOffset := len(g.root) + 1

	relPathFn := func(path string) (relPath string) {
		if len(path) >= relPathOffset {
			relPath = path[relPathOffset:]
		}
		return
	}

	idx, err := g.repo.Storer.Index()
	if err != nil {
		return fmt.Errorf("failed to open git index: %w", err)
	}

	for path := range g.paths {

		err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(g.root, path)
			if err != nil {
				return err
			}

			if _, err = idx.Entry(relPath); errors.Is(err, index.ErrEntryNotFound) {
				// we skip this path as it's not staged
				log.Debugf("Path not found in git index, skipping: %v, %v", relPath, path)
				return nil
			}

			file := File{
				Path:    path,
				RelPath: relPathFn(path),
				Info:    info,
			}

			return fn(&file, err)
		})
		if err != nil {
			return err
		}

	}

	return nil
}

func NewGit(root string, paths chan string) (Walker, error) {
	repo, err := git.PlainOpen(root)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repo: %w", err)
	}
	return &gitWalker{root, paths, repo}, nil
}
