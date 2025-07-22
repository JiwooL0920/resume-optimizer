import { useCallback } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { AppDispatch, RootState } from '../store'
import { 
  fetchApiKeys, 
  createApiKey, 
  deleteApiKey 
} from '../store/slices/apiKeysSlice'
import { 
  fetchResumes, 
  uploadResume, 
  deleteResume, 
  selectResume 
} from '../store/slices/resumeSlice'
import { optimizeResume } from '../store/slices/optimizationSlice'
import { ApiKey, Resume, OptimizationRequest } from '../types'

// Custom hook for API key operations
export const useApiKeys = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { keys, isLoading, error } = useSelector((state: RootState) => state.apiKeys)

  const fetchKeys = useCallback(() => {
    dispatch(fetchApiKeys())
  }, [dispatch])

  const addKey = useCallback((provider: string, apiKey: string) => {
    return dispatch(createApiKey({ provider, api_key: apiKey }))
  }, [dispatch])

  const deleteKey = useCallback((keyId: string) => {
    return dispatch(deleteApiKey(keyId))
  }, [dispatch])

  const getKeysByProvider = useCallback((provider: string) => {
    return keys.filter(key => key.provider === provider)
  }, [keys])

  const hasKeyForProvider = useCallback((provider: string) => {
    return keys.some(key => key.provider === provider)
  }, [keys])

  return {
    keys,
    isLoading,
    error,
    fetchKeys,
    addKey,
    deleteKey,
    getKeysByProvider,
    hasKeyForProvider,
  }
}

// Custom hook for resume operations
export const useResumes = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { resumes, selectedResume, isLoading, error, uploadProgress } = useSelector(
    (state: RootState) => state.resume
  )

  const fetchResumeList = useCallback(() => {
    dispatch(fetchResumes())
  }, [dispatch])

  const uploadNewResume = useCallback((file: File) => {
    return dispatch(uploadResume(file))
  }, [dispatch])

  const selectResumeById = useCallback((resume: Resume) => {
    dispatch(selectResume(resume))
  }, [dispatch])

  const deleteResumeById = useCallback((resumeId: string) => {
    return dispatch(deleteResume(resumeId))
  }, [dispatch])

  const getResumeById = useCallback((resumeId: string) => {
    return resumes.find(resume => resume.id === resumeId)
  }, [resumes])

  return {
    resumes,
    selectedResume,
    uploadProgress,
    isLoading,
    error,
    fetchResumeList,
    uploadNewResume,
    selectResumeById,
    deleteResumeById,
    getResumeById,
  }
}

// Custom hook for optimization operations
export const useOptimization = () => {
  const dispatch = useDispatch<AppDispatch>()
  const optimizationState = useSelector((state: RootState) => state.optimization)

  const startOptimization = useCallback((request: OptimizationRequest) => {
    return dispatch(optimizeResume(request))
  }, [dispatch])

  return {
    startOptimization,
    ...optimizationState,
  }
}

// Custom hook for authentication operations
export const useAuth = () => {
  const authState = useSelector((state: RootState) => state.auth)
  
  const isAuthenticated = authState.isAuthenticated
  const user = authState.user
  const token = authState.token
  const isLoading = authState.isLoading
  const error = authState.error

  return {
    isAuthenticated,
    user,
    token,
    isLoading,
    error,
  }
}

// Custom hook for error handling
export const useError = () => {
  const handleApiError = useCallback((error: any): string => {
    if (error?.response?.data?.error) {
      return error.response.data.error
    }
    if (error?.message) {
      return error.message
    }
    if (typeof error === 'string') {
      return error
    }
    return 'An unexpected error occurred'
  }, [])

  const isNetworkError = useCallback((error: any): boolean => {
    return error?.code === 'NETWORK_ERROR' || 
           error?.message?.includes('Network Error') ||
           !navigator.onLine
  }, [])

  return {
    handleApiError,
    isNetworkError,
  }
}

// Custom hook for localStorage operations
export const useLocalStorage = () => {
  const setItem = useCallback((key: string, value: any) => {
    try {
      localStorage.setItem(key, JSON.stringify(value))
    } catch (error) {
      console.error('Failed to save to localStorage:', error)
    }
  }, [])

  const getItem = useCallback((key: string) => {
    try {
      const item = localStorage.getItem(key)
      return item ? JSON.parse(item) : null
    } catch (error) {
      console.error('Failed to read from localStorage:', error)
      return null
    }
  }, [])

  const removeItem = useCallback((key: string) => {
    try {
      localStorage.removeItem(key)
    } catch (error) {
      console.error('Failed to remove from localStorage:', error)
    }
  }, [])

  return {
    setItem,
    getItem,
    removeItem,
  }
}

// Custom hook for form validation
export const useFormValidation = () => {
  const validateEmail = useCallback((email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
  }, [])

  const validateUrl = useCallback((url: string): boolean => {
    try {
      new URL(url)
      return true
    } catch {
      return false
    }
  }, [])

  const validateApiKey = useCallback((key: string, provider: string): boolean => {
    if (!key || key.trim().length === 0) return false
    
    switch (provider) {
      case 'openai':
        return key.startsWith('sk-')
      case 'anthropic':
        return key.startsWith('sk-ant-')
      default:
        return key.length > 10 // Basic length check
    }
  }, [])

  return {
    validateEmail,
    validateUrl,
    validateApiKey,
  }
}