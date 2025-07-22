import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { setToken, fetchProfile } from '../../store/slices/authSlice'
import { AppDispatch, RootState } from '../../store'
import { useNavigate } from 'react-router-dom'

const AuthCallback: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const { isLoading, user, error } = useSelector((state: RootState) => state.auth)
  const [tokenProcessed, setTokenProcessed] = useState(false)

  useEffect(() => {
    // Extract token from URL query parameters
    const urlParams = new URLSearchParams(window.location.search)
    const token = urlParams.get('token')
    
    if (token && !tokenProcessed) {
      setTokenProcessed(true)
      
      // Store the token and mark as authenticated
      dispatch(setToken(token))
      
      // Fetch user profile data using the token directly
      dispatch(fetchProfile(token))
      
      // Clear the URL parameters
      window.history.replaceState({}, document.title, window.location.pathname)
    } else if (!token) {
      // If no token, something went wrong with auth
      console.error('No token received from auth callback')
      navigate('/', { replace: true })
    }
  }, [dispatch, tokenProcessed, navigate])

  // Redirect to dashboard when profile is loaded successfully
  useEffect(() => {
    if (user && !isLoading && tokenProcessed) {
      navigate('/dashboard', { replace: true })
    } else if (error && tokenProcessed) {
      console.error('Authentication failed:', error)
      navigate('/', { replace: true })
    }
  }, [user, isLoading, error, tokenProcessed, navigate])

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p className="mt-4 text-gray-600">Completing authentication...</p>
      </div>
    </div>
  )
}

export default AuthCallback