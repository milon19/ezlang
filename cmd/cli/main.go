package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/milon19/ezlang/internal/application/services"
	"github.com/milon19/ezlang/internal/config"
	"github.com/milon19/ezlang/pkg/translation"
)

func main() {
	configPath := flag.String("config", ".ezlang.yml", "Path to configuration file")
	rewriteMain := flag.Bool("rewrite", false, "Rewrite main file. Default is false")
	flag.Parse()

	cfg, err := config.Load(*configPath)

	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	ctx := context.Background()
	tc, err := translation.NewTranslationClient(ctx)
	if err != nil {
		fmt.Println("Error creating translation client")
		return
	}

	for _, file := range cfg.Files {
		err := services.ProcessFile(file, tc)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", file.Path, err)
		} else {
			fmt.Printf("Successfully processed file %s\n", file.Path)
		}
		if rewriteMain != nil && *rewriteMain {
			err := services.RewriteMainFile(file)
			if err != nil {
				fmt.Printf("Error rewriting main file: %v\n", err)
			} else {
				fmt.Printf("Successfully rewrote main file\n")
			}
		}
	}

}
