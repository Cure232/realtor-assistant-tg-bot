package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.sajari.com/docconv"
)

// Supported file extensions for docconv
var supportedExtensions = map[string]bool{
	".pdf":  true,
	".docx": true,
	".doc":  true,
	".rtf":  true,
	".odt":  true,
	".html": true,
}

func main() {
	// Specify the folder containing documents
	inputDir := "./documents"
	outputDir := "./output"

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Walk through the input directory
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return nil // Continue walking despite the error
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file extension is supported
		ext := strings.ToLower(filepath.Ext(path))
		if !supportedExtensions[ext] {
			log.Printf("Skipping unsupported file: %s", path)
			return nil
		}

		// Parse the document
		if err := parseDocument(path, outputDir); err != nil {
			log.Printf("Failed to parse %s: %v", path, err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through directory: %v", err)
	}

	fmt.Println("Document parsing completed.")
}

// parseDocument processes a single document and saves the extracted text
func parseDocument(filePath, outputDir string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Convert document to text using docconv TODO OTHER FORMATS AND META??
	res, _, err := docconv.ConvertPDF(file)
	if err != nil {
		return fmt.Errorf("failed to convert document: %w", err)
	}

	// Generate output file path
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	outputPath := filepath.Join(outputDir, baseName+".txt")

	// Write extracted text to output file
	if err := ioutil.WriteFile(outputPath, []byte(res), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Printf("Successfully parsed %s, output saved to %s", filePath, outputPath)
	return nil
}
