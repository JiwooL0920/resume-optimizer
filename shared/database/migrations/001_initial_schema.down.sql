-- Rollback initial schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_user_api_keys_updated_at ON user_api_keys;
DROP TRIGGER IF EXISTS update_optimization_sessions_updated_at ON optimization_sessions;
DROP TRIGGER IF EXISTS update_resumes_updated_at ON resumes;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_user_api_keys_provider;
DROP INDEX IF EXISTS idx_user_api_keys_user_id;
DROP INDEX IF EXISTS idx_feedback_session_id;
DROP INDEX IF EXISTS idx_optimization_sessions_created_at;
DROP INDEX IF EXISTS idx_optimization_sessions_status;
DROP INDEX IF EXISTS idx_optimization_sessions_resume_id;
DROP INDEX IF EXISTS idx_optimization_sessions_user_id;
DROP INDEX IF EXISTS idx_resumes_created_at;
DROP INDEX IF EXISTS idx_resumes_user_id;
DROP INDEX IF EXISTS idx_users_google_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS feedback;
DROP TABLE IF EXISTS user_api_keys;
DROP TABLE IF EXISTS optimization_sessions;
DROP TABLE IF EXISTS resumes;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";