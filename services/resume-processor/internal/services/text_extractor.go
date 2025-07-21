package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

// TextExtractor handles text extraction from various file formats
type TextExtractor struct{}

// NewTextExtractor creates a new TextExtractor instance
func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

// ExtractText extracts plain text from a file based on its extension
func (te *TextExtractor) ExtractText(filePath string) (string, error) {
	fmt.Printf("=== NEW TEXT EXTRACTOR CALLED ===\n")
	fmt.Printf("File: %s\n", filePath)
	ext := strings.ToLower(filepath.Ext(filePath))
	fmt.Printf("Extension detected: %s\n", ext)
	
	switch ext {
	case ".pdf":
		fmt.Printf("Calling UNIPDF extraction...\n")
		return te.extractFromPDF(filePath)
	case ".txt":
		return te.extractFromText(filePath)
	default:
		// For unknown formats, try to read as text and clean up
		content, err := te.extractFromText(filePath)
		if err != nil {
			return "", fmt.Errorf("unsupported file format: %s", ext)
		}
		return te.cleanTextContent(content), nil
	}
}

// extractFromPDF extracts text from PDF files using improved filtering
func (te *TextExtractor) extractFromPDF(filePath string) (string, error) {
	fmt.Printf("=== STARTING IMPROVED PDF EXTRACTION ===\n")
	fmt.Printf("Processing file: %s\n", filePath)
	
	// Open PDF file
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		fmt.Printf("ERROR: Failed to open PDF: %v\n", err)
		return "", fmt.Errorf("failed to open PDF: %v", err)
	}
	defer file.Close()
	
	fmt.Printf("PDF opened successfully, pages: %d\n", reader.NumPage())
	
	var textBuilder strings.Builder
	validTextFound := false
	
	// Extract text from all pages
	for pageNum := 1; pageNum <= reader.NumPage(); pageNum++ {
		fmt.Printf("Processing page %d/%d\n", pageNum, reader.NumPage())
		
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			fmt.Printf("Page %d: Null page, skipping\n", pageNum)
			continue
		}
		
		// Get page content
		content := page.Content()
		if content.Text == nil {
			fmt.Printf("Page %d: No text objects found\n", pageNum)
			continue
		}
		
		pageTextExtracted := 0
		totalSegments := 0
		validSegments := 0
		// Extract text from the page with advanced filtering
		for _, text := range content.Text {
			if text.S == "" {
				continue
			}
			
			totalSegments++
			fmt.Printf("  Segment %d: '%s' (len=%d)\n", totalSegments, text.S, len(text.S))
			
			// Apply intelligent filtering
			if te.isValidTextSegment(text.S) {
				validSegments++
				cleanedSegment := te.cleanTextSegment(text.S)
				if cleanedSegment != "" {
					fmt.Printf("  -> Valid segment %d: '%s'\n", validSegments, cleanedSegment)
					textBuilder.WriteString(cleanedSegment + " ")
					pageTextExtracted += len(cleanedSegment)
					validTextFound = true
				} else {
					fmt.Printf("  -> Segment valid but cleaned to empty\n")
				}
			} else {
				fmt.Printf("  -> Segment filtered out by validation\n")
			}
		}
		
		if pageTextExtracted > 0 {
			fmt.Printf("Page %d extracted: %d characters from %d/%d segments\n", pageNum, pageTextExtracted, validSegments, totalSegments)
			textBuilder.WriteString("\n") // Add line break between pages
		} else {
			fmt.Printf("Page %d: No valid text content found (processed %d segments, %d valid)\n", pageNum, totalSegments, validSegments)
		}
	}
	
	rawText := textBuilder.String()
	cleanedText := te.cleanTextContent(rawText)
	
	// If no valid text found, try more aggressive extraction methods
	if !validTextFound || len(cleanedText) < 10 {
		fmt.Printf("No text found via PDF parsing, trying aggressive extraction...\n")
		
		// Try aggressive text extraction first
		aggressiveText, err := te.extractTextAggressive(filePath)
		if err == nil && len(aggressiveText) >= 10 {
			fmt.Printf("Aggressive extraction successful: %d characters\n", len(aggressiveText))
			cleanedText = te.cleanTextContent(aggressiveText)
			validTextFound = true
		} else {
			fmt.Printf("Aggressive extraction failed or insufficient text, attempting OCR...\n")
			// Try OCR as last resort
			ocrText, err := te.extractTextWithOCR(filePath)
			if err != nil {
				fmt.Printf("OCR extraction also failed: %v\n", err)
				return "", fmt.Errorf("PDF text extraction failed - no readable text found via PDF parsing (%d characters), aggressive extraction, or OCR. This PDF might be image-based, corrupted, or have very complex formatting", len(cleanedText))
			}
			
			if len(ocrText) < 10 {
				return "", fmt.Errorf("PDF text extraction failed - OCR found only %d characters. This PDF might be corrupted or contain no readable text", len(ocrText))
			}
			
			fmt.Printf("OCR extraction successful: %d characters\n", len(ocrText))
			cleanedText = te.cleanTextContent(ocrText)
			validTextFound = true
		}
	}
	
	// Additional validation - check for meaningful content
	words := strings.Fields(cleanedText)
	if len(words) < 5 {
		return "", fmt.Errorf("PDF text extraction failed - insufficient text content (only %d words). This PDF might be image-based", len(words))
	}
	
	// Check if text looks like actual resume content vs PDF artifacts
	if te.looksLikePDFArtifacts(cleanedText) {
		return "", fmt.Errorf("PDF text extraction failed - extracted content appears to be PDF structure rather than readable text")
	}
	
	// Debug logging
	fmt.Printf("=== PDF EXTRACTION RESULTS ===\n")
	fmt.Printf("Total pages processed: %d\n", reader.NumPage())
	fmt.Printf("Raw text length: %d characters\n", len(rawText))
	fmt.Printf("Cleaned text length: %d characters\n", len(cleanedText))
	fmt.Printf("Word count: %d words\n", len(words))
	fmt.Printf("Valid text found: %t\n", validTextFound)
	
	// Write extracted text to debug file
	debugFile := "/tmp/extracted_text_debug.txt"
	if err := os.WriteFile(debugFile, []byte(cleanedText), 0644); err != nil {
		fmt.Printf("Warning: Could not write debug file: %v\n", err)
	} else {
		fmt.Printf("Debug: Extracted text written to %s\n", debugFile)
	}
	
	// Show text sample for debugging
	if len(cleanedText) > 500 {
		fmt.Printf("Text sample (first 500 chars):\n%s...\n", cleanedText[:500])
	} else {
		fmt.Printf("Complete extracted text:\n%s\n", cleanedText)
	}
	fmt.Printf("=== END EXTRACTION ===\n")
	
	return cleanedText, nil
}

// extractFromText reads plain text files
func (te *TextExtractor) extractFromText(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %v", err)
	}
	
	return string(content), nil
}

// cleanTextContent cleans up extracted text content for better readability
func (te *TextExtractor) cleanTextContent(text string) string {
	fmt.Printf("=== STARTING TEXT CLEANING ===\n")
	fmt.Printf("Original text length: %d characters\n", len(text))
	
	// Step 1: Remove weird symbols and artifacts
	text = te.removeArtifacts(text)
	
	// Step 2: Fix spacing issues
	text = te.fixSpacing(text)
	
	// Step 3: Normalize line breaks and structure
	text = te.normalizeStructure(text)
	
	// Step 4: Final cleanup
	text = te.finalCleanup(text)
	
	fmt.Printf("Cleaned text length: %d characters\n", len(text))
	fmt.Printf("=== TEXT CLEANING COMPLETE ===\n")
	
	return text
}

// removeArtifacts removes PDF artifacts and weird symbols
func (te *TextExtractor) removeArtifacts(text string) string {
	// Remove common PDF symbols and artifacts
	artifacts := map[string]string{
		"Ω": " ",     // Replace Ω with space
		"⊗": "•",     // Replace ⊗ with bullet
		"∗": "•",     // Replace ∗ with bullet  
		"§": "",      // Remove section symbol
		"ï": "",      // Remove weird i
		"#": "",      // Remove hash in context
		"&": "and",   // Replace & with and
		"R&D": "R&D", // Keep R&D as is
	}
	
	for artifact, replacement := range artifacts {
		if artifact == "R&D" {
			continue // Skip R&D replacement
		}
		text = strings.ReplaceAll(text, artifact, replacement)
	}
	
	// Handle R&D specially
	text = strings.ReplaceAll(text, "R and D", "R&D")
	
	return text
}

// fixSpacing fixes spacing issues between characters
func (te *TextExtractor) fixSpacing(text string) string {
	fmt.Printf("Fixing spacing issues...\n")
	originalLen := len(text)
	
	// Step 1: Fix extremely spaced out text (every character separated)
	// This handles cases like "J i w o o L e e" -> "JiwooLee"
	extremeSpacingPattern := regexp.MustCompile(`\b([A-Za-z])(?:\s+([A-Za-z]))+\b`)
	text = extremeSpacingPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Remove ALL spaces within the match
		return strings.ReplaceAll(match, " ", "")
	})
	
	// Step 2: Fix spaced digits in phone numbers and years
	digitSpacingPattern := regexp.MustCompile(`\b(\d)(?:\s+(\d))+\b`)
	text = digitSpacingPattern.ReplaceAllStringFunc(text, func(match string) string {
		return strings.ReplaceAll(match, " ", "")
	})
	
	// Step 3: Fix email addresses specifically  
	// Handle patterns like "j i w o o @ g m a i l . c o m"
	emailPattern := regexp.MustCompile(`([a-zA-Z0-9]+(?:\s+[a-zA-Z0-9]+)*)\s*@\s*([a-zA-Z]+(?:\s+[a-zA-Z]+)*)\s*\.\s*([a-zA-Z]+(?:\s+[a-zA-Z]+)*)`)
	text = emailPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Remove spaces in email components
		parts := strings.Split(match, "@")
		if len(parts) == 2 {
			localPart := strings.ReplaceAll(parts[0], " ", "")
			domainPart := strings.ReplaceAll(parts[1], " ", "")
			return localPart + "@" + domainPart
		}
		return strings.ReplaceAll(match, " ", "")
	})
	
	// Step 4: Fix URLs like "l i n k e d i n . c o m"
	urlPattern := regexp.MustCompile(`\b([a-zA-Z]+(?:\s+[a-zA-Z]+)*)\s*\.\s*([a-zA-Z]+(?:\s+[a-zA-Z]+)*)(?:\s*\.\s*([a-zA-Z]+(?:\s+[a-zA-Z]+)*))*(?:\s*/\s*([a-zA-Z0-9]+(?:\s+[a-zA-Z0-9]+)*))*`)
	text = urlPattern.ReplaceAllStringFunc(text, func(match string) string {
		return strings.ReplaceAll(match, " ", "")
	})
	
	// Step 5: Fix specific common words that get spaced out
	commonWords := map[string]string{
		"S u m m a r y":     "Summary",
		"E x p e r i e n c e": "Experience", 
		"E d u c a t i o n":   "Education",
		"S k i l l s":       "Skills",
		"P r o j e c t s":    "Projects",
	}
	
	for spaced, fixed := range commonWords {
		text = strings.ReplaceAll(text, spaced, fixed)
	}
	
	// Step 6: More aggressive general spacing fix
	// Fix any remaining patterns of single characters separated by spaces
	singleCharPattern := regexp.MustCompile(`\b([A-Za-z])\s+([A-Za-z])\s+([A-Za-z])\s+([A-Za-z])\s*([A-Za-z]*)\b`)
	text = singleCharPattern.ReplaceAllStringFunc(text, func(match string) string {
		return strings.ReplaceAll(match, " ", "")
	})
	
	// Step 7: Fix remaining 3-character spaced patterns
	threeCharPattern := regexp.MustCompile(`\b([A-Za-z])\s+([A-Za-z])\s+([A-Za-z])\b`)
	text = threeCharPattern.ReplaceAllStringFunc(text, func(match string) string {
		return strings.ReplaceAll(match, " ", "")
	})
	
	fmt.Printf("Spacing fix: %d -> %d characters\n", originalLen, len(text))
	return text
}

// normalizeStructure improves text structure and formatting
func (te *TextExtractor) normalizeStructure(text string) string {
	// Split into lines and clean each line
	lines := strings.Split(text, "\n")
	var cleanLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Add proper spacing after sections
		if te.isSectionHeader(line) {
			if len(cleanLines) > 0 {
				cleanLines = append(cleanLines, "") // Add blank line before section
			}
			cleanLines = append(cleanLines, line)
			cleanLines = append(cleanLines, "") // Add blank line after section
		} else {
			cleanLines = append(cleanLines, line)
		}
	}
	
	return strings.Join(cleanLines, "\n")
}

// isSectionHeader checks if a line is likely a section header
func (te *TextExtractor) isSectionHeader(line string) bool {
	sectionHeaders := []string{
		"Summary", "Technical Skills", "Experience", "Education", 
		"Certifications", "Projects", "Skills", "Work Experience",
		"Professional Experience", "Leadership", "Activity",
	}
	
	for _, header := range sectionHeaders {
		if strings.Contains(strings.ToLower(line), strings.ToLower(header)) {
			return true
		}
	}
	return false
}

// finalCleanup performs final text cleanup
func (te *TextExtractor) finalCleanup(text string) string {
	// Remove control characters and non-printable characters
	text = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return -1 // Remove character
		}
		return r
	}, text)
	
	// Clean up excessive whitespace but preserve structure
	lines := strings.Split(text, "\n")
	var finalLines []string
	
	for _, line := range lines {
		// Clean up spaces within lines but keep single spaces
		line = regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")
		line = strings.TrimSpace(line)
		finalLines = append(finalLines, line)
	}
	
	// Remove excessive empty lines (max 2 consecutive)
	result := strings.Join(finalLines, "\n")
	excessiveNewlines := regexp.MustCompile(`\n{3,}`)
	result = excessiveNewlines.ReplaceAllString(result, "\n\n")
	
	return strings.TrimSpace(result)
}

// ValidateTextLength checks if extracted text is reasonable for AI processing
func (te *TextExtractor) ValidateTextLength(text string) error {
	fmt.Printf("=== TEXT VALIDATION DEBUG ===\n")
	fmt.Printf("Text length: %d characters\n", len(text))
	
	if len(text) == 0 {
		fmt.Printf("ERROR: No text content found\n")
		return fmt.Errorf("no text content found in the file")
	}
	
	if len(text) < 50 {
		fmt.Printf("ERROR: Text too short\n")
		return fmt.Errorf("file content is too short to be a valid resume")
	}
	
	// Temporarily increase limit to 300,000 characters for debugging
	if len(text) > 300000 {
		fmt.Printf("ERROR: Text too large (%d characters)\n", len(text))
		return fmt.Errorf("file content is too large (over 300,000 characters)")
	}
	
	fmt.Printf("Text validation passed\n")
	return nil
}

// isValidTextSegment checks if a text segment is likely to be actual readable text
func (te *TextExtractor) isValidTextSegment(text string) bool {
	// Skip empty or very short segments
	if len(strings.TrimSpace(text)) < 2 {
		return false
	}
	
	// Skip common PDF artifacts
	pdfArtifacts := []string{
		"obj", "endobj", "stream", "endstream", "xref", "trailer",
		"%%PDF", "/Type", "/Font", "/Page", "/Catalog", "/Length",
		"/Filter", "/FlateDecode", "/ASCIIHexDecode", "/ASCII85Decode",
		"BT", "ET", "Tf", "Td", "TJ", "Tj", "cm", "q", "Q",
		"<<", ">>", "null", "true", "false", "R",
	}
	
	for _, artifact := range pdfArtifacts {
		if strings.Contains(text, artifact) {
			return false
		}
	}
	
	// Skip if mostly non-printable characters
	printableCount := 0
	for _, r := range text {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			printableCount++
		}
	}
	if float64(printableCount)/float64(len(text)) < 0.7 {
		return false
	}
	
	// Skip if looks like hex or encoded data
	hexPattern := regexp.MustCompile(`^[0-9a-fA-F\s]+$`)
	if len(text) > 10 && hexPattern.MatchString(strings.ReplaceAll(text, " ", "")) {
		return false
	}
	
	// Skip if mostly numbers and special characters (likely coordinates or encoding)
	wordPattern := regexp.MustCompile(`[a-zA-Z]{2,}`)
	if !wordPattern.MatchString(text) && len(text) > 5 {
		return false
	}
	
	return true
}

// cleanTextSegment cleans a valid text segment
func (te *TextExtractor) cleanTextSegment(text string) string {
	// Remove extra whitespace
	cleaned := strings.TrimSpace(text)
	
	// Remove common PDF operators that might slip through
	pdfOperators := []string{"BT", "ET", "Tf", "Td", "TJ", "Tj"}
	for _, op := range pdfOperators {
		cleaned = strings.ReplaceAll(cleaned, op, "")
	}
	
	// Clean up extra spaces
	spacePattern := regexp.MustCompile(`\s+`)
	cleaned = spacePattern.ReplaceAllString(cleaned, " ")
	
	return strings.TrimSpace(cleaned)
}

// looksLikePDFArtifacts checks if the final extracted text still looks like PDF structure
func (te *TextExtractor) looksLikePDFArtifacts(text string) bool {
	// Check for high concentration of PDF-like patterns
	pdfPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\d+\s+\d+\s+obj`),
		regexp.MustCompile(`<</[A-Za-z]+`),
		regexp.MustCompile(`\[\d+\s+\d+\s+R\]`),
		regexp.MustCompile(`/[A-Z][A-Za-z]+\s+\d+`),
	}
	
	matches := 0
	for _, pattern := range pdfPatterns {
		if pattern.MatchString(text) {
			matches++
		}
	}
	
	// If we find multiple PDF patterns, it's likely still PDF structure
	if matches >= 2 {
		return true
	}
	
	// Check if text has very low ratio of actual words to total length
	words := strings.Fields(text)
	wordChars := 0
	for _, word := range words {
		// Count characters in words that look like actual words (contain letters)
		if regexp.MustCompile(`[a-zA-Z]`).MatchString(word) {
			wordChars += len(word)
		}
	}
	
	// If less than 30% of characters are in actual words, it's probably artifacts
	if len(text) > 100 && float64(wordChars)/float64(len(text)) < 0.3 {
		return true
	}
	
	return false
}

// extractTextAggressive extracts all available text without filtering (for difficult PDFs)
func (te *TextExtractor) extractTextAggressive(filePath string) (string, error) {
	fmt.Printf("=== STARTING AGGRESSIVE EXTRACTION ===\n")
	
	// Open PDF file
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %v", err)
	}
	defer file.Close()
	
	fmt.Printf("Aggressive extraction: PDF has %d pages\n", reader.NumPage())
	
	var textBuilder strings.Builder
	allSegments := 0
	
	// Extract ALL text segments without filtering
	for pageNum := 1; pageNum <= reader.NumPage(); pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}
		
		content := page.Content()
		if content.Text == nil {
			continue
		}
		
		pageSegments := 0
		for _, text := range content.Text {
			if text.S != "" {
				allSegments++
				pageSegments++
				
				// Add ALL text segments, even ones that might look like artifacts
				// We'll clean them up later
				textBuilder.WriteString(text.S + " ")
			}
		}
		
		fmt.Printf("Page %d: extracted %d segments\n", pageNum, pageSegments)
		if pageSegments > 0 {
			textBuilder.WriteString("\n")
		}
	}
	
	rawText := textBuilder.String()
	fmt.Printf("Aggressive extraction complete: %d segments, %d characters\n", allSegments, len(rawText))
	
	// Apply basic cleaning but less aggressive filtering
	cleanedText := te.cleanTextContentAggressive(rawText)
	
	fmt.Printf("After aggressive cleaning: %d characters\n", len(cleanedText))
	
	if len(cleanedText) > 100 {
		previewLen := 500
		if len(cleanedText) < previewLen {
			previewLen = len(cleanedText)
		}
		fmt.Printf("Aggressive sample (first %d chars):\n%s...\n", previewLen, cleanedText[:previewLen])
	} else {
		fmt.Printf("Complete aggressive text:\n%s\n", cleanedText)
	}
	
	return cleanedText, nil
}

// cleanTextContentAggressive applies minimal cleaning for aggressive extraction
func (te *TextExtractor) cleanTextContentAggressive(text string) string {
	// Remove obvious PDF artifacts but keep most text
	pdfJunk := []string{
		"obj", "endobj", "stream", "endstream",
		"<<", ">>", "/Type", "/Page", "/Font",
		"BT", "ET", "Tf", "Td",
	}
	
	for _, junk := range pdfJunk {
		text = strings.ReplaceAll(text, junk, " ")
	}
	
	// Clean up excessive whitespace
	spacePattern := regexp.MustCompile(`\s+`)
	text = spacePattern.ReplaceAllString(text, " ")
	
	// Remove control characters but be less aggressive
	result := strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return ' ' // Replace with space instead of removing
		}
		return r
	}, text)
	
	return strings.TrimSpace(result)
}

// extractTextWithOCR extracts text from PDF using OCR (for image-based PDFs)
func (te *TextExtractor) extractTextWithOCR(filePath string) (string, error) {
	fmt.Printf("=== STARTING OCR EXTRACTION ===\n")
	
	// Create temporary directory for OCR processing
	tempDir := "/tmp/ocr_" + filepath.Base(filePath)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Step 1: Convert PDF to images using pdfimages
	fmt.Printf("Converting PDF to images...\n")
	imagePrefix := filepath.Join(tempDir, "page")
	cmd := exec.Command("pdfimages", "-png", filePath, imagePrefix)
	if err := cmd.Run(); err != nil {
		// Try alternative: convert PDF using ImageMagick
		fmt.Printf("pdfimages failed, trying ImageMagick convert...\n")
		imageFile := filepath.Join(tempDir, "page.png")
		cmd = exec.Command("convert", "-density", "300", filePath, imageFile)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to convert PDF to images: %v", err)
		}
	}
	
	// Step 2: Find generated image files
	files, err := filepath.Glob(filepath.Join(tempDir, "*.png"))
	if err != nil {
		return "", fmt.Errorf("failed to find image files: %v", err)
	}
	if len(files) == 0 {
		// Try jpg files as fallback
		files, _ = filepath.Glob(filepath.Join(tempDir, "*.jpg"))
		if len(files) == 0 {
			return "", fmt.Errorf("no image files generated from PDF")
		}
	}
	
	fmt.Printf("Found %d image files to process\n", len(files))
	
	// Step 3: OCR each image file
	var allText strings.Builder
	for i, imageFile := range files {
		fmt.Printf("Processing image %d/%d: %s\n", i+1, len(files), imageFile)
		
		// Use tesseract to extract text
		cmd = exec.Command("tesseract", imageFile, "stdout", "-c", "preserve_interword_spaces=1")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Warning: OCR failed for %s: %v\n", imageFile, err)
			continue
		}
		
		text := strings.TrimSpace(string(output))
		if text != "" {
			allText.WriteString(text)
			allText.WriteString("\n\n")
			fmt.Printf("Extracted %d characters from image %d\n", len(text), i+1)
		} else {
			fmt.Printf("No text found in image %d\n", i+1)
		}
	}
	
	ocrResult := allText.String()
	fmt.Printf("=== OCR EXTRACTION COMPLETE ===\n")
	fmt.Printf("Total OCR text extracted: %d characters\n", len(ocrResult))
	
	if len(ocrResult) > 500 {
		fmt.Printf("OCR sample (first 500 chars):\n%s...\n", ocrResult[:500])
	} else {
		fmt.Printf("Complete OCR text:\n%s\n", ocrResult)
	}
	
	return ocrResult, nil
}
