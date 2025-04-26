package services

import (
	"fmt"
	"github.com/milon19/ezlang/internal/config"
	"github.com/milon19/ezlang/pkg/po"
	"github.com/milon19/ezlang/pkg/translation"
	"os"
	"strings"
)

func ProcessFile(file config.FileConfig, tc *translation.Client) error {
	fmt.Println("Processing file:", file.Path)

	poFile, err := po.ReadPoFile(file.Path)
	if err != nil {
		return fmt.Errorf("error reading PO file: %w", err)
	}

	outputPath := strings.Split(file.Path, ".")[0] + "_output" + ".po"

	if err := po.WritePOFile(outputPath, poFile, file.Lang, tc); err != nil {
		return fmt.Errorf("error writing PO file: %w", err)
	}
	return nil
}

func RewriteMainFile(file config.FileConfig) error {
	err := os.Remove(file.Path)
	if err != nil {
		return err
	}

	outputPath := strings.Split(file.Path, ".")[0] + "_output" + ".po"
	err = os.Rename(outputPath, file.Path)
	if err != nil {
		return err
	}
	return nil
}
