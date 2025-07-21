import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'

interface User {
  id: string
  email: string
  name: string
  picture_url?: string
}

interface AuthState {
  user: User | null
  token: string | null
  isLoading: boolean
  error: string | null
  isAuthenticated: boolean
}

const initialState: AuthState = {
  user: null,
  token: null, // Clear any stored token for clean start
  isLoading: false,
  error: null,
  isAuthenticated: false,
}

export const loginWithGoogle = createAsyncThunk(
  'auth/loginWithGoogle',
  async () => {
    window.location.href = 'http://localhost:8080/api/v1/auth/google'
  }
)

export const logout = createAsyncThunk(
  'auth/logout',
  async () => {
    const response = await fetch('http://localhost:8080/api/v1/auth/logout', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (response.ok) {
      localStorage.removeItem('token')
      return true
    }
    throw new Error('Logout failed')
  }
)

export const fetchProfile = createAsyncThunk(
  'auth/fetchProfile',
  async (token?: string) => {
    const authToken = token || localStorage.getItem('token')
    const response = await fetch('http://localhost:8080/api/v1/auth/profile', {
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    })
    
    if (response.ok) {
      const data = await response.json()
      return data.user
    }
    throw new Error('Failed to fetch profile')
  }
)

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload
      state.isAuthenticated = true
      localStorage.setItem('token', action.payload)
    },
    clearAuth: (state) => {
      state.user = null
      state.token = null
      state.isAuthenticated = false
      localStorage.removeItem('token')
    },
    clearError: (state) => {
      state.error = null
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchProfile.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(fetchProfile.fulfilled, (state, action) => {
        state.isLoading = false
        state.user = action.payload
        state.isAuthenticated = true
      })
      .addCase(fetchProfile.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.error.message || 'Failed to fetch profile'
        state.isAuthenticated = false
      })
      .addCase(logout.fulfilled, (state) => {
        state.user = null
        state.token = null
        state.isAuthenticated = false
      })
  },
})

export const { setToken, clearAuth, clearError } = authSlice.actions
export default authSlice.reducer