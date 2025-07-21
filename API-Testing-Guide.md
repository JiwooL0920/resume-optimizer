# Resume Optimizer API Testing Guide

This guide will help you set up and test the Resume Optimizer APIs using Postman.

## Setup Instructions

### 1. Import Collection and Environment

1. **Import the Collection:**
   - Open Postman
   - Click "Import" in the top left
   - Select `Resume-Optimizer-API.postman_collection.json`

2. **Import the Environment:**
   - Click "Import" again
   - Select `Resume-Optimizer.postman_environment.json`

3. **Select Environment:**
   - In the top right dropdown, select "Resume Optimizer Environment"

### 2. Start Your Services

Make sure your services are running:
```bash
./start-services.sh
```

This should start:
- Auth Service on port 8080
- Resume Processor Service on port 8081

## Testing Workflow

### Step 1: Basic Health Checks

Test that your services are running:

1. **Auth Service Health Check**
   - Send GET request to `{{auth_service_url}}/health`
   - Expected: 404 (no health endpoint implemented yet, but service is responding)

2. **Resume Service Health Check** 
   - Send GET request to `{{resume_service_url}}/health`
   - Expected: 404 (no health endpoint implemented yet, but service is responding)

### Step 2: Test Resume Management

Start with the resume processor service since it doesn't require authentication:

1. **List All Resumes**
   - Send GET request to `{{resume_service_url}}/api/v1/resumes/`
   - Expected: `{"resumes": []}` (empty array if no resumes uploaded yet)

2. **Upload a Resume**
   - Use POST request to `{{resume_service_url}}/api/v1/resumes/upload`
   - In Body tab, select "form-data"
   - Add key "file" with type "File" and select a resume file (PDF, DOC, DOCX)
   - Expected response: `{"id": "some-uuid"}`
   - **Copy the ID from the response and update the `resume_id` environment variable**

3. **Get Resume by ID**
   - Send GET request to `{{resume_service_url}}/api/v1/resumes/{{resume_id}}`
   - Expected: Resume details including file info

4. **List Resumes Again**
   - Send GET request to `{{resume_service_url}}/api/v1/resumes/`
   - Expected: Array containing your uploaded resume

### Step 3: Test Optimization Endpoints

1. **Optimize Resume**
   - Send POST request to `{{resume_service_url}}/api/v1/optimize/`
   - Expected: `{"message": "Resume optimization in development"}` (placeholder response)

2. **Apply Feedback**
   - Send POST request to `{{resume_service_url}}/api/v1/optimize/feedback`
   - Expected: `{"message": "Apply feedback in development"}` (placeholder response)

### Step 4: Authentication Testing (Advanced)

⚠️ **Note:** Google OAuth requires proper configuration and browser interaction.

For authentication testing, you'll need:

1. **Configure Google OAuth:**
   - Update your `.env` file with valid Google Client ID and Secret
   - Set proper redirect URL

2. **Browser-based OAuth Flow:**
   - Open `http://localhost:8080/auth/google` in a browser
   - Complete Google login
   - The callback will contain the JWT token
   - Copy the token and update the `jwt_token` environment variable

3. **Test Protected Endpoint:**
   - Send GET request to `{{auth_service_url}}/profile`
   - Include Authorization header with your JWT token
   - Expected: User profile information

## Environment Variables

Update these variables in your Postman environment as you test:

- `auth_service_url`: http://localhost:8080
- `resume_service_url`: http://localhost:8081  
- `jwt_token`: JWT token from Google OAuth callback
- `resume_id`: ID returned from resume upload

## Common Issues & Solutions

### Service Not Responding
- Check that both services are running with `./start-services.sh`
- Verify PostgreSQL is running and accessible
- Check terminal for error messages

### File Upload Issues
- Ensure you're using "form-data" body type
- File key should be exactly "file"
- Supported formats: PDF, DOC, DOCX

### Database Connection Issues
- Verify PostgreSQL port-forward is active: `lsof -i :5432`
- Check database credentials in `.env` file
- Ensure database is accessible from localhost:5432

### Authentication Issues
- Google OAuth requires valid client credentials
- Redirect URL must match Google OAuth configuration
- JWT tokens have expiration times

## API Endpoints Summary

### Auth Service (Port 8080)
- `GET /auth/google` - Initiate Google OAuth
- `GET /profile` - Get user profile (requires JWT)
- `POST /logout` - Logout user

### Resume Processor Service (Port 8081)
- `GET /api/v1/resumes/` - List all resumes
- `POST /api/v1/resumes/upload` - Upload resume file
- `GET /api/v1/resumes/:id` - Get resume by ID
- `DELETE /api/v1/resumes/:id` - Delete resume by ID
- `POST /api/v1/optimize/` - Optimize resume (development)
- `POST /api/v1/optimize/feedback` - Apply feedback (development)

## Next Steps

Once basic functionality is working:
1. Implement health check endpoints
2. Add authentication middleware to resume endpoints
3. Develop AI optimization features
4. Add error handling and validation
5. Implement rate limiting and security measures

Happy testing! 🚀
