package pkg

import (
	"fmt"
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
	content, err := os.ReadFile(j.OriginPath)
	if err != nil {
		return fmt.Errorf("error trying to read %s: %w", j.OriginPath, err)
	}

	if err := os.WriteFile(j.FinalPath, content, j.Mode); err != nil {
		return fmt.Errorf("error writing %s: %w", j.FinalPath, err)
	}

	return nil
}

type Copier struct {
	TargetDir       string
	IgnoreFilePaths []string
	wg              sync.WaitGroup
}

func NewCopier(targetDir string, ignoreFilePaths []string) *Copier {
	return &Copier{
		TargetDir:       targetDir,
		IgnoreFilePaths: ignoreFilePaths,
	}
}

func (c *Copier) ListDir(rootDir string, jobChan chan<- CopyJob) {
	defer c.wg.Done()

	dirs, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", rootDir, err)
		return
	}

	for _, dir := range dirs {
		fullOriginPath := filepath.Join(rootDir, dir.Name())
		relPath, err := filepath.Rel(rootDir, fullOriginPath)
		if err != nil {
			fmt.Printf("Error trying to calculate relative path: %v\n", err)
			return
		}

		fullFinalPath := filepath.Join(c.TargetDir, relPath)

		if c.ignoreFilePath(dir.Name()) {
			continue
		}

		if dir.IsDir() {
			info, err := dir.Info()
			if err == nil {
				os.MkdirAll(fullFinalPath, info.Mode())
			}

			c.wg.Add(1)
			go c.ListDir(fullOriginPath, jobChan)
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
	ignoreList := append(c.IgnoreFilePaths, ".DS_Store", "Thumbs.db", ".git", "node_modules", "dist", "build")
	return slices.Contains(ignoreList, filePath)
}

func (c *Copier) Start(rootDir string, jobChan chan CopyJob) {
	c.wg.Add(1)
	go c.ListDir(rootDir, jobChan)

	go func() {
		c.wg.Wait()
		close(jobChan)
	}()
}
