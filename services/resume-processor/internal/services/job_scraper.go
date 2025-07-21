package services

import (
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// JobScraper handles fetching and extracting job descriptions from URLs
type JobScraper struct {
	client *http.Client
}

// NewJobScraper creates a new JobScraper instance
func NewJobScraper() *JobScraper {
	return &JobScraper{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchJobDescription fetches and extracts job description from a URL
func (js *JobScraper) FetchJobDescription(url string) (string, error) {
	// Set user agent to avoid blocking
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := js.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse HTML and extract meaningful text content
	content, err := js.extractTextFromHTML(string(body))
	if err != nil {
		return "", err
	}

	// Clean and filter the content to focus on job description
	cleanContent := js.cleanJobDescription(content)
	
	return cleanContent, nil
}

// extractTextFromHTML extracts text content from HTML
func (js *JobScraper) extractTextFromHTML(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var content strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				content.WriteString(text + " ")
			}
		}
		// Skip script and style tags
		if n.Data != "script" && n.Data != "style" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)

	return content.String(), nil
}

// cleanJobDescription cleans and filters the extracted content
func (js *JobScraper) cleanJobDescription(content string) string {
	// Remove extra whitespace
	content = strings.Join(strings.Fields(content), " ")
	
	// Common job description keywords to help identify relevant content
	jobKeywords := []string{
		"responsibilities", "requirements", "qualifications", "experience", 
		"skills", "duties", "role", "position", "job description", 
		"what you'll do", "about the role", "we are looking for",
		"candidate", "preferred", "required", "minimum", "years",
	}
	
	// Find sections that contain job-related keywords
	words := strings.Fields(strings.ToLower(content))
	var relevantSections []string
	var currentSection []string
	
	for i, word := range words {
		// Check if current word or nearby words contain job keywords
		isRelevant := false
		for _, keyword := range jobKeywords {
			if strings.Contains(word, keyword) {
				isRelevant = true
				break
			}
		}
		
		if isRelevant {
			// Start a new relevant section
			if len(currentSection) > 0 {
				relevantSections = append(relevantSections, strings.Join(currentSection, " "))
			}
			currentSection = []string{}
			
			// Include surrounding context (50 words before and after)
			start := max(0, i-50)
			end := min(len(words), i+100)
			currentSection = append(currentSection, words[start:end]...)
		}
	}
	
	if len(currentSection) > 0 {
		relevantSections = append(relevantSections, strings.Join(currentSection, " "))
	}
	
	if len(relevantSections) > 0 {
		result := strings.Join(relevantSections, "\n\n")
		// Limit to reasonable length (3000 characters)
		if len(result) > 3000 {
			result = result[:3000] + "..."
		}
		return result
	}
	
	// Fallback: return first 2000 characters if no specific sections found
	if len(content) > 2000 {
		content = content[:2000] + "..."
	}
	
	return content
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}