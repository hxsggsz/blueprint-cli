package main

import (
	"blueprint/config"
	"blueprint/pkg"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		return
	}
	bluePrintConfig, err := cfg.GetConfig()
	if err != nil {
		fmt.Printf("Error getting config: %v\n", err)
		return
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error trying to get current directory: %v\n", err)
		return
	}
	c := pkg.NewCopier(currentDir, bluePrintConfig.IgnoreFilePaths)
	templates := c.ListTemplates(bluePrintConfig.TemplatePath)
	selectedTemplate := selectTemplate(templates)

	jobChan := make(chan pkg.CopyJob, 100)

	c.Start(filepath.Join(bluePrintConfig.TemplatePath, selectedTemplate), jobChan)

	numWorkers := runtime.NumCPU() * 2
	var workerWg sync.WaitGroup

	start := time.Now()
	for range numWorkers {
		workerWg.Go(func() {
			for job := range jobChan {
				if err := job.CreateFile(); err != nil {
					fmt.Printf("Error trying to read %s: %v\n", job.OriginPath, err)
				}
			}
		})
	}

	workerWg.Wait()
	duration := time.Since(start)
	fmt.Printf("Done in: %v\n", duration.Round(time.Millisecond))
}

func selectTemplate(templates []string) string {
	fmt.Println("Available templates:")
	for i, template := range templates {
		fmt.Printf("%d: %s\n", i+1, template)
	}

	var choice int
	fmt.Print("Select a template by number: ")
	_, err := fmt.Scan(&choice)
	if err != nil || choice < 1 || choice > len(templates) {
		fmt.Println("Invalid selection. Please try again.")
		return selectTemplate(templates)
	}

	return templates[choice-1]
}
