package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AIOptimizer handles AI-based resume optimization
type AIOptimizer struct {
	client *http.Client
}

// NewAIOptimizer creates a new AIOptimizer instance
func NewAIOptimizer() *AIOptimizer {
	return &AIOptimizer{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// OptimizationRequest represents a request to optimize a resume
type OptimizationRequest struct {
	ResumeContent      string `json:"resume_content"`
	JobDescription     string `json:"job_description"`
	AIModel           string `json:"ai_model"`
	KeepOnePage       bool   `json:"keep_one_page"`
	UserAPIKey        string `json:"user_api_key"`
}

// OptimizationResponse represents the response from AI optimization
type OptimizationResponse struct {
	OptimizedContent string `json:"optimized_content"`
	Summary         string `json:"summary"`
	Changes         []string `json:"changes"`
}

// OptimizeResume optimizes a resume using the specified AI model
func (ai *AIOptimizer) OptimizeResume(req OptimizationRequest) (*OptimizationResponse, error) {
	switch {
	case strings.HasPrefix(req.AIModel, "gpt-"):
		return ai.optimizeWithOpenAI(req)
	case strings.HasPrefix(req.AIModel, "claude-"):
		return ai.optimizeWithClaude(req)
	default:
		return nil, fmt.Errorf("unsupported AI model: %s", req.AIModel)
	}
}

// optimizeWithOpenAI handles optimization using OpenAI GPT models
func (ai *AIOptimizer) optimizeWithOpenAI(req OptimizationRequest) (*OptimizationResponse, error) {
	prompt := ai.buildOptimizationPrompt(req.ResumeContent, req.JobDescription, req.KeepOnePage)
	
	requestBody := map[string]interface{}{
		"model": req.AIModel,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "You are an expert resume writer and career coach. Your task is to optimize resumes to better match job descriptions while maintaining authenticity and improving the candidate's chances of getting noticed by HR and ATS systems.",
			},
			{
				"role": "user",
				"content": prompt,
			},
		},
		"max_tokens": 4000,
		"temperature": 0.7,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+req.UserAPIKey)

	resp, err := ai.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return ai.parseOptimizationResponse(openAIResp.Choices[0].Message.Content)
}

// optimizeWithClaude handles optimization using Anthropic Claude models
func (ai *AIOptimizer) optimizeWithClaude(req OptimizationRequest) (*OptimizationResponse, error) {
	prompt := ai.buildOptimizationPrompt(req.ResumeContent, req.JobDescription, req.KeepOnePage)
	
	requestBody := map[string]interface{}{
		"model": req.AIModel,
		"max_tokens": 4000,
		"temperature": 0.7,
		"messages": []map[string]string{
			{
				"role": "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", req.UserAPIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := ai.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var claudeResp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, err
	}

	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("no response from Claude")
	}

	return ai.parseOptimizationResponse(claudeResp.Content[0].Text)
}

// buildOptimizationPrompt creates the prompt for AI optimization
func (ai *AIOptimizer) buildOptimizationPrompt(resumeContent, jobDescription string, keepOnePage bool) string {
	pageLimitText := ""
	if keepOnePage {
		pageLimitText = "\n- IMPORTANT: Keep the optimized resume to exactly ONE PAGE. Be selective and concise."
	}

	return fmt.Sprintf(`You are an expert resume writer and career coach. Please optimize the following resume to better match the given job description while maintaining authenticity and improving ATS (Applicant Tracking System) compatibility.

CURRENT RESUME:
%s

JOB DESCRIPTION:
%s

OPTIMIZATION REQUIREMENTS:
- Tailor the resume to highlight relevant skills and experiences for this specific job
- Use keywords from the job description naturally throughout the resume
- Improve the format and structure for better ATS compatibility
- Enhance action verbs and quantify achievements where possible
- Maintain truthfulness - do not add fake experiences or skills%s

Please provide your response in the following JSON format:
{
  "optimized_content": "The complete optimized resume content in clean text format",
  "summary": "A brief summary of the main changes made and why they improve the candidate's chances",
  "changes": ["List of specific changes made", "Each change as a separate item", "Focus on the most impactful modifications"]
}

Ensure the optimized_content is ready to be used as-is and maintains professional formatting.`, resumeContent, jobDescription, pageLimitText)
}

// parseOptimizationResponse parses the AI response and extracts the structured data
func (ai *AIOptimizer) parseOptimizationResponse(content string) (*OptimizationResponse, error) {
	// Try to extract JSON from the response
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	
	if startIdx == -1 || endIdx == -1 {
		// Fallback: treat entire response as optimized content
		return &OptimizationResponse{
			OptimizedContent: content,
			Summary:         "Resume optimized successfully",
			Changes:         []string{"Resume has been tailored to match the job requirements"},
		}, nil
	}

	jsonStr := content[startIdx:endIdx+1]
	
	var response OptimizationResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		// Fallback: treat entire response as optimized content
		return &OptimizationResponse{
			OptimizedContent: content,
			Summary:         "Resume optimized successfully",
			Changes:         []string{"Resume has been tailored to match the job requirements"},
		}, nil
	}

	return &response, nil
}