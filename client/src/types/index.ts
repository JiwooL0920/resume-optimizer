// Shared type definitions for the Resume Optimizer application

export interface User {
  id: string;
  email: string;
  google_id?: string;
  name: string;
  picture_url?: string;
  created_at: string;
  updated_at: string;
}

export interface Resume {
  id: string;
  user_id: string;
  title: string;
  original_content: string;
  extracted_text: string;
  file_type: string;
  file_size?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApiKey {
  id: string;
  provider: string;
  masked_key: string;
  created_at: string;
}

export interface OptimizationSession {
  id: string;
  user_id: string;
  resume_id: string;
  job_description_url?: string;
  job_description_text?: string;
  ai_model: string;
  keep_one_page: boolean;
  optimized_content?: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  created_at: string;
  updated_at: string;
}

export interface OptimizationResult {
  session: OptimizationSession;
  summary?: string;
  changes?: string[];
}

export interface Feedback {
  id: string;
  session_id: string;
  section_highlight: string;
  user_comment: string;
  is_processed: boolean;
  created_at: string;
}

// API Response types
export interface AuthResponse {
  user: User;
  token: string;
}

export interface ResumeListResponse {
  resumes: Resume[];
}

export interface ApiKeyListResponse {
  api_keys: ApiKey[];
}

// Form types
export interface LoginFormData {
  email: string;
  password: string;
}

export interface ResumeUploadData {
  file: File;
}

export interface OptimizationRequest {
  resumeId: string;
  jobDescriptionUrl?: string;
  jobDescriptionText?: string;
  aiModel: string;
  keepOnePage: boolean;
  userApiKey: string;
}

export interface FeedbackRequest {
  sessionId: string;
  sectionHighlight: string;
  userComment: string;
}

// UI State types
export interface LoadingState {
  isLoading: boolean;
  error?: string;
}

export interface ErrorState {
  message: string;
  code?: string;
  timestamp: string;
}

// AI Model options
export type AIModel = 'gpt-4' | 'gpt-3.5-turbo' | 'claude-3-opus' | 'claude-3-sonnet';

export interface AIModelOption {
  value: AIModel;
  label: string;
  provider: string;
}

// Navigation types
export type PageRoute = 'dashboard' | 'optimize' | 'settings';