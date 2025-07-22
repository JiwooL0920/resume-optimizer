// Centralized API configuration

// Base URLs for microservices
export const API_ENDPOINTS = {
  AUTH_BASE: process.env.REACT_APP_AUTH_API_URL || 'http://localhost:8080/api/v1',
  RESUME_BASE: process.env.REACT_APP_RESUME_API_URL || 'http://localhost:8081/api/v1',
};

// Auth service endpoints
export const AUTH_ENDPOINTS = {
  LOGIN: `${API_ENDPOINTS.AUTH_BASE}/auth/login`,
  LOGOUT: `${API_ENDPOINTS.AUTH_BASE}/auth/logout`,
  GOOGLE_LOGIN: `${API_ENDPOINTS.AUTH_BASE}/auth/google`,
  CALLBACK: `${API_ENDPOINTS.AUTH_BASE}/auth/callback`,
  PROFILE: `${API_ENDPOINTS.AUTH_BASE}/auth/profile`,
  API_KEYS: `${API_ENDPOINTS.AUTH_BASE}/api-keys`,
};

// Resume service endpoints
export const RESUME_ENDPOINTS = {
  RESUMES: `${API_ENDPOINTS.RESUME_BASE}/resumes`,
  UPLOAD: `${API_ENDPOINTS.RESUME_BASE}/resumes/upload`,
  OPTIMIZE: `${API_ENDPOINTS.RESUME_BASE}/resumes/optimize`,
  FEEDBACK: `${API_ENDPOINTS.RESUME_BASE}/feedback`,
};

// HTTP headers
export const DEFAULT_HEADERS = {
  'Content-Type': 'application/json',
};

// Request configuration
export const REQUEST_CONFIG = {
  timeout: 30000, // 30 seconds
  retryAttempts: 3,
};

// Helper function to get authorization headers
export const getAuthHeaders = (token?: string) => {
  const authToken = token || localStorage.getItem('token');
  return authToken 
    ? { ...DEFAULT_HEADERS, Authorization: `Bearer ${authToken}` }
    : DEFAULT_HEADERS;
};

// Helper function to build full endpoint URL
export const buildEndpoint = (path: string, baseUrl = API_ENDPOINTS.RESUME_BASE) => {
  return `${baseUrl}${path}`;
};

// Error status codes
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
} as const;

// AI Model configurations
export const AI_MODELS = {
  'gpt-4': {
    label: 'GPT-4',
    provider: 'openai',
    requiresApiKey: true,
  },
  'gpt-3.5-turbo': {
    label: 'GPT-3.5 Turbo',
    provider: 'openai',
    requiresApiKey: true,
  },
  'claude-3-opus': {
    label: 'Claude 3 Opus',
    provider: 'anthropic',
    requiresApiKey: true,
  },
  'claude-3-sonnet': {
    label: 'Claude 3 Sonnet',
    provider: 'anthropic',
    requiresApiKey: true,
  },
} as const;

export type AIModelKey = keyof typeof AI_MODELS;