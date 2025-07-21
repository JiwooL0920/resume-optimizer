import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'

interface OptimizationSession {
  id: string
  user_id: string
  resume_id: string
  job_description_url?: string
  job_description_text?: string
  ai_model: string
  keep_one_page: boolean
  optimized_content?: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  created_at: string
  updated_at: string
}

interface Feedback {
  id: string
  session_id: string
  section_highlight: string
  user_comment: string
  is_processed: boolean
  created_at: string
}

interface OptimizationState {
  currentSession: OptimizationSession | null
  sessions: OptimizationSession[]
  feedback: Feedback[]
  isOptimizing: boolean
  isProcessingFeedback: boolean
  error: string | null
  selectedAiModel: string
  keepOnePage: boolean
}

const initialState: OptimizationState = {
  currentSession: null,
  sessions: [],
  feedback: [],
  isOptimizing: false,
  isProcessingFeedback: false,
  error: null,
  selectedAiModel: 'gpt-4',
  keepOnePage: false,
}

export const optimizeResume = createAsyncThunk(
  'optimization/optimizeResume',
  async (params: {
    resumeId: string
    jobDescriptionUrl?: string
    jobDescriptionText?: string
    aiModel: string
    keepOnePage: boolean
  }) => {
    const response = await fetch('/api/v1/optimize/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(params)
    })
    
    if (response.ok) {
      return await response.json()
    }
    throw new Error('Failed to optimize resume')
  }
)

export const applyFeedback = createAsyncThunk(
  'optimization/applyFeedback',
  async (params: {
    sessionId: string
    feedback: Array<{
      section_highlight: string
      user_comment: string
    }>
  }) => {
    const response = await fetch('/api/v1/optimize/feedback', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(params)
    })
    
    if (response.ok) {
      return await response.json()
    }
    throw new Error('Failed to apply feedback')
  }
)

const optimizationSlice = createSlice({
  name: 'optimization',
  initialState,
  reducers: {
    setCurrentSession: (state, action: PayloadAction<OptimizationSession>) => {
      state.currentSession = action.payload
    },
    clearCurrentSession: (state) => {
      state.currentSession = null
      state.feedback = []
    },
    setSelectedAiModel: (state, action: PayloadAction<string>) => {
      state.selectedAiModel = action.payload
    },
    setKeepOnePage: (state, action: PayloadAction<boolean>) => {
      state.keepOnePage = action.payload
    },
    addFeedback: (state, action: PayloadAction<Feedback>) => {
      state.feedback.push(action.payload)
    },
    removeFeedback: (state, action: PayloadAction<string>) => {
      state.feedback = state.feedback.filter(f => f.id !== action.payload)
    },
    clearError: (state) => {
      state.error = null
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(optimizeResume.pending, (state) => {
        state.isOptimizing = true
        state.error = null
      })
      .addCase(optimizeResume.fulfilled, (state, action) => {
        state.isOptimizing = false
        state.currentSession = action.payload
        state.sessions.push(action.payload)
      })
      .addCase(optimizeResume.rejected, (state, action) => {
        state.isOptimizing = false
        state.error = action.error.message || 'Failed to optimize resume'
      })
      .addCase(applyFeedback.pending, (state) => {
        state.isProcessingFeedback = true
        state.error = null
      })
      .addCase(applyFeedback.fulfilled, (state, action) => {
        state.isProcessingFeedback = false
        if (state.currentSession) {
          state.currentSession = { ...state.currentSession, ...action.payload }
        }
        state.feedback = []
      })
      .addCase(applyFeedback.rejected, (state, action) => {
        state.isProcessingFeedback = false
        state.error = action.error.message || 'Failed to apply feedback'
      })
  },
})

export const {
  setCurrentSession,
  clearCurrentSession,
  setSelectedAiModel,
  setKeepOnePage,
  addFeedback,
  removeFeedback,
  clearError
} = optimizationSlice.actions

export default optimizationSlice.reducer