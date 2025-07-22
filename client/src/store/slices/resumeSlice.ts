import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'
import { Resume } from '../../types'

interface ResumeState {
  resumes: Resume[]
  selectedResume: Resume | null
  isLoading: boolean
  error: string | null
  uploadProgress: number
}

const initialState: ResumeState = {
  resumes: [],
  selectedResume: null,
  isLoading: false,
  error: null,
  uploadProgress: 0,
}

export const fetchResumes = createAsyncThunk(
  'resume/fetchResumes',
  async () => {
    const response = await fetch('http://localhost:8081/api/v1/resumes/', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (response.ok) {
      return await response.json()
    }
    throw new Error('Failed to fetch resumes')
  }
)

export const uploadResume = createAsyncThunk(
  'resume/uploadResume',
  async (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await fetch('http://localhost:8081/api/v1/resumes/upload', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: formData
    })
    
    if (response.ok) {
      return await response.json()
    }
    throw new Error('Failed to upload resume')
  }
)

export const deleteResume = createAsyncThunk(
  'resume/deleteResume',
  async (resumeId: string) => {
    const response = await fetch(`http://localhost:8081/api/v1/resumes/${resumeId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (response.ok) {
      return resumeId
    }
    throw new Error('Failed to delete resume')
  }
)

const resumeSlice = createSlice({
  name: 'resume',
  initialState,
  reducers: {
    selectResume: (state, action: PayloadAction<Resume>) => {
      state.selectedResume = action.payload
    },
    clearSelectedResume: (state) => {
      state.selectedResume = null
    },
    setUploadProgress: (state, action: PayloadAction<number>) => {
      state.uploadProgress = action.payload
    },
    clearError: (state) => {
      state.error = null
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchResumes.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(fetchResumes.fulfilled, (state, action) => {
        state.isLoading = false
        state.resumes = action.payload.resumes || []
      })
      .addCase(fetchResumes.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.error.message || 'Failed to fetch resumes'
      })
      .addCase(uploadResume.pending, (state) => {
        state.isLoading = true
        state.error = null
        state.uploadProgress = 0
      })
      .addCase(uploadResume.fulfilled, (state, action) => {
        state.isLoading = false
        state.resumes.push(action.payload)
        state.uploadProgress = 100
      })
      .addCase(uploadResume.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.error.message || 'Failed to upload resume'
        state.uploadProgress = 0
      })
      .addCase(deleteResume.fulfilled, (state, action) => {
        state.resumes = state.resumes.filter(resume => resume.id !== action.payload)
        if (state.selectedResume?.id === action.payload) {
          state.selectedResume = null
        }
      })
  },
})

export const { selectResume, clearSelectedResume, setUploadProgress, clearError } = resumeSlice.actions
export default resumeSlice.reducer