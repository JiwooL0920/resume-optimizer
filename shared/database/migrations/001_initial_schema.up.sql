-- Initial schema for resume optimizer application

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    google_id VARCHAR(255) UNIQUE,
    name VARCHAR(255) NOT NULL,
    picture_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Resumes table
CREATE TABLE resumes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    original_content TEXT NOT NULL,
    extracted_text TEXT,
    file_type VARCHAR(10) NOT NULL DEFAULT '.pdf',
    file_size INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Optimization sessions table
CREATE TABLE optimization_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    job_description_url VARCHAR(1000),
    job_description_text TEXT,
    ai_model VARCHAR(50) NOT NULL,
    keep_one_page BOOLEAN DEFAULT false,
    optimized_content TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Feedback table
CREATE TABLE feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES optimization_sessions(id) ON DELETE CASCADE,
    section_highlight TEXT NOT NULL,
    user_comment TEXT NOT NULL,
    is_processed BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- API keys table (encrypted)
CREATE TABLE user_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    encrypted_key TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, provider)
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_resumes_user_id ON resumes(user_id);
CREATE INDEX idx_resumes_created_at ON resumes(created_at);
CREATE INDEX idx_optimization_sessions_user_id ON optimization_sessions(user_id);
CREATE INDEX idx_optimization_sessions_resume_id ON optimization_sessions(resume_id);
CREATE INDEX idx_optimization_sessions_status ON optimization_sessions(status);
CREATE INDEX idx_optimization_sessions_created_at ON optimization_sessions(created_at);
CREATE INDEX idx_feedback_session_id ON feedback(session_id);
CREATE INDEX idx_user_api_keys_user_id ON user_api_keys(user_id);
CREATE INDEX idx_user_api_keys_provider ON user_api_keys(provider);

-- Add constraints
ALTER TABLE optimization_sessions ADD CONSTRAINT check_status 
    CHECK (status IN ('pending', 'processing', 'completed', 'failed'));

ALTER TABLE optimization_sessions ADD CONSTRAINT check_job_description 
    CHECK (job_description_url IS NOT NULL OR job_description_text IS NOT NULL);

ALTER TABLE user_api_keys ADD CONSTRAINT check_provider 
    CHECK (provider IN ('openai', 'anthropic', 'google', 'cohere'));

-- Update timestamps trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_resumes_updated_at BEFORE UPDATE ON resumes 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_optimization_sessions_updated_at BEFORE UPDATE ON optimization_sessions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_api_keys_updated_at BEFORE UPDATE ON user_api_keys 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();