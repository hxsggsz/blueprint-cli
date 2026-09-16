package main

import (
	"blueprint/config"
	"blueprint/pkg"
	"fmt"
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
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Printf("Error trying to get current directory: %v\n", err)
	// 	return
	// }
	c := pkg.NewCopier("destino", bluePrintConfig.IgnoreFilePaths)
	jobChan := make(chan pkg.CopyJob, 100)

	c.Start(bluePrintConfig.TemplatePath, jobChan)

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
	fmt.Printf("Done in: %v\n", duration)
}
