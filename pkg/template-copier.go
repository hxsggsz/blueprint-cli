package pkg

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

type CopyJob struct {
	OriginPath string
	FinalPath  string
	Mode       fs.FileMode
}

func (j CopyJob) CreateFile() error {
	src, err := os.Open(j.OriginPath)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", j.OriginPath, err)
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(j.FinalPath), 0755); err != nil {
		return fmt.Errorf("error creating parent directory for %s: %w", j.FinalPath, err)
	}

	dst, err := os.OpenFile(j.FinalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, j.Mode)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %w", j.FinalPath, err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error during streaming to %s: %w", j.FinalPath, err)
	}

	return nil
}

type Copier struct {
	TargetDir       string
	IgnoreFilePaths []string
	wg              sync.WaitGroup
}

func NewCopier(targetDir string, ignoreFilePaths []string) *Copier {
	defaultIgnores := []string{".DS_Store", "Thumbs.db", ".git", "node_modules", "dist", "build"}
	allIgnores := append(ignoreFilePaths, defaultIgnores...)

	return &Copier{
		TargetDir:       targetDir,
		IgnoreFilePaths: allIgnores,
	}
}
func (c *Copier) ListTemplates(currentDir string) []string {
	templates := make([]string, 0)
	dirs, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", currentDir, err)
		return []string{}
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			templates = append(templates, dir.Name())
		}
	}

	return templates
}

func (c *Copier) ListDir(baseDir, currentDir string, jobChan chan<- CopyJob) {
	defer c.wg.Done()

	dirs, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", currentDir, err)
		return
	}

	for _, dir := range dirs {
		if c.ignoreFilePath(dir.Name()) {
			continue
		}

		fullOriginPath := filepath.Join(currentDir, dir.Name())

		relPath, err := filepath.Rel(baseDir, fullOriginPath)
		if err != nil {
			fmt.Printf("Error trying to calculate relative path: %v\n", err)
			return
		}

		fullFinalPath := filepath.Join(c.TargetDir, relPath)

		if dir.IsDir() {
			info, err := dir.Info()
			if err == nil {
				os.MkdirAll(fullFinalPath, info.Mode())
			}

			c.wg.Add(1)
			go c.ListDir(baseDir, fullOriginPath, jobChan)
		} else {
			info, err := dir.Info()
			if err != nil {
				fmt.Printf("Error getting information of %s: %v\n", dir.Name(), err)
				continue
			}

			jobChan <- CopyJob{
				OriginPath: fullOriginPath,
				FinalPath:  fullFinalPath,
				Mode:       info.Mode(),
			}
		}
	}
}

func (c *Copier) ignoreFilePath(filePath string) bool {
	return slices.Contains(c.IgnoreFilePaths, filePath)
}

func (c *Copier) Start(rootDir string, jobChan chan CopyJob) {
	c.wg.Add(1)
	go c.ListDir(rootDir, rootDir, jobChan)

	go func() {
		c.wg.Wait()
		close(jobChan)
	}()
}
