import React, { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { setToken, fetchProfile } from '../../store/slices/authSlice'
import { AppDispatch } from '../../store'

const AuthCallback: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()

  useEffect(() => {
    // Extract token from URL query parameters
    const urlParams = new URLSearchParams(window.location.search)
    const token = urlParams.get('token')
    
    if (token) {
      // Store the token and mark as authenticated
      dispatch(setToken(token))
      
      // Fetch user profile data using the token directly
      dispatch(fetchProfile(token))
      
      // Clear the URL parameters and redirect to dashboard
      window.history.replaceState({}, document.title, '/')
    } else {
      // If no token, something went wrong with auth
      console.error('No token received from auth callback')
      // Redirect to login or show error
      window.location.href = '/'
    }
  }, [dispatch])

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