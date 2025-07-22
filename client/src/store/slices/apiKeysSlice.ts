import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit'

interface ApiKey {
  id: string
  provider: string
  masked_key: string
  created_at: string
}

interface ApiKeysState {
  keys: ApiKey[]
  isLoading: boolean
  error: string | null
}

const initialState: ApiKeysState = {
  keys: [],
  isLoading: false,
  error: null,
}

export const fetchApiKeys = createAsyncThunk(
  'apiKeys/fetchApiKeys',
  async (_, { rejectWithValue }) => {
    try {
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/user/api-keys', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        return data.api_keys || []
      } else {
        throw new Error('Failed to fetch API keys')
      }
    } catch (error: any) {
      return rejectWithValue(error.message)
    }
  }
)

export const createApiKey = createAsyncThunk(
  'apiKeys/createApiKey',
  async ({ provider, api_key }: { provider: string, api_key: string }, { rejectWithValue }) => {
    try {
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/user/api-keys', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          provider,
          api_key
        })
      })

      if (response.ok) {
        const data = await response.json()
        return data.api_key
      } else {
        const error = await response.json()
        throw new Error(error.error || 'Failed to add API key')
      }
    } catch (error: any) {
      return rejectWithValue(error.message)
    }
  }
)

export const deleteApiKey = createAsyncThunk(
  'apiKeys/deleteApiKey',
  async (keyId: string, { rejectWithValue }) => {
    try {
      const token = localStorage.getItem('token')
      const response = await fetch(`http://localhost:8080/api/v1/user/api-keys/${keyId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (response.ok) {
        return keyId
      } else {
        throw new Error('Failed to delete API key')
      }
    } catch (error: any) {
      return rejectWithValue(error.message)
    }
  }
)

const apiKeysSlice = createSlice({
  name: 'apiKeys',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null
    },
    clearApiKeys: (state) => {
      state.keys = []
      state.error = null
      state.isLoading = false
    }
  },
  extraReducers: (builder) => {
    builder
      // Fetch API keys
      .addCase(fetchApiKeys.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(fetchApiKeys.fulfilled, (state, action) => {
        state.isLoading = false
        state.keys = action.payload
      })
      .addCase(fetchApiKeys.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.payload as string
      })
      // Create API key
      .addCase(createApiKey.pending, (state) => {
        state.isLoading = true
        state.error = null
      })
      .addCase(createApiKey.fulfilled, (state, action) => {
        state.isLoading = false
        state.keys.push(action.payload)
      })
      .addCase(createApiKey.rejected, (state, action) => {
        state.isLoading = false
        state.error = action.payload as string
      })
      // Delete API key
      .addCase(deleteApiKey.pending, (state) => {
        state.error = null
      })
      .addCase(deleteApiKey.fulfilled, (state, action) => {
        state.keys = state.keys.filter(key => key.id !== action.payload)
      })
      .addCase(deleteApiKey.rejected, (state, action) => {
        state.error = action.payload as string
      })
  },
})

export const { clearError, clearApiKeys } = apiKeysSlice.actions
export default apiKeysSlice.reducer
