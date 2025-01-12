package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/phcarneirobc/profile-search/internal/config"
	"github.com/phcarneirobc/profile-search/internal/model"
	"github.com/phcarneirobc/profile-search/internal/platform"
)

func main() {
	username := flag.String("username", "", "Username to search")
	flag.Parse()

	if *username == "" {
		fmt.Println("Por favor, forneça um username usando -username")
		return
	}

	cfg := config.GetDefaultConfig()
	platforms := platform.GetPlatforms()

	results := make(chan model.Result, len(platforms))
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, cfg.ConcurrentRequests)

	fmt.Printf("\nProcurando por '%s' em %d plataformas...\n\n", *username, len(platforms))

	for _, p := range platforms {
		wg.Add(1)
		go func(p platform.Platform) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			platform.CheckPlatform(p, *username, cfg, results)
		}(p)
	}

	var allResults []model.Result
	go func() {
		for result := range results {
			allResults = append(allResults, result)
			if result.Error != "" {
				fmt.Printf("❌ %s: %s (%.2fs)\n", result.Platform, result.Error, result.ResponseTime.Seconds())
			} else if result.Exists {
				fmt.Printf("✅ %s: Perfil encontrado - %s (%.2fs)\n", result.Platform, result.URL, result.ResponseTime.Seconds())
				if result.Info != nil {
					if result.Info.Name != "" {
						fmt.Printf("   Nome: %s\n", result.Info.Name)
					}
					if result.Info.Location != "" {
						fmt.Printf("   Localização: %s\n", result.Info.Location)
					}
					if result.Info.Bio != "" {
						fmt.Printf("   Bio: %s\n", result.Info.Bio)
					}
				}
			} else {
				fmt.Printf("❌ %s: Perfil não encontrado (%.2fs)\n", result.Platform, result.ResponseTime.Seconds())
			}
		}
	}()

	wg.Wait()
	close(results)

	var totalTime time.Duration
	found := 0
	for _, r := range allResults {
		totalTime += r.ResponseTime
		if r.Exists {
			found++
		}
	}

	fmt.Printf("\nResultados:\n")
	fmt.Printf("Total de plataformas verificadas: %d\n", len(allResults))
	fmt.Printf("Perfis encontrados: %d\n", found)
	fmt.Printf("Tempo médio por requisição: %.2fs\n", totalTime.Seconds()/float64(len(allResults)))
}
