import React, { useEffect } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { Routes, Route, Navigate, useLocation } from 'react-router-dom'
import { RootState } from '../../store'
import { fetchProfile } from '../../store/slices/authSlice'
import { AppDispatch } from '../../store'
import Header from './Header'
import Login from '../Auth/Login'
import AuthCallback from '../Auth/AuthCallback'
import Dashboard from '../Dashboard/Dashboard'
import OptimizeResume from '../Optimize/OptimizeResume'
import Settings from '../Settings/Settings'

// Protected Route component
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useSelector((state: RootState) => state.auth)
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

const Layout: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { isAuthenticated, token } = useSelector((state: RootState) => state.auth)
  const location = useLocation()

  useEffect(() => {
    if (token && !isAuthenticated) {
      dispatch(fetchProfile())
    }
  }, [dispatch, token, isAuthenticated])

  // Get current page from pathname for header highlighting
  const getCurrentPage = () => {
    const path = location.pathname
    if (path.includes('/optimize')) return 'optimize'
    if (path.includes('/settings')) return 'settings'
    return 'dashboard'
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {isAuthenticated && <Header currentPage={getCurrentPage()} />}
      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <Routes>
          {/* Public routes */}
          <Route path="/login" element={<Login />} />
          <Route path="/auth/callback" element={<AuthCallback />} />
          
          {/* Protected routes */}
          <Route path="/dashboard" element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          } />
          <Route path="/optimize" element={
            <ProtectedRoute>
              <OptimizeResume />
            </ProtectedRoute>
          } />
          <Route path="/settings" element={
            <ProtectedRoute>
              <Settings />
            </ProtectedRoute>
          } />
          
          {/* Default redirects */}
          <Route path="/" element={
            isAuthenticated ? <Navigate to="/dashboard" replace /> : <Navigate to="/login" replace />
          } />
          <Route path="*" element={
            isAuthenticated ? <Navigate to="/dashboard" replace /> : <Navigate to="/login" replace />
          } />
        </Routes>
      </main>
    </div>
  )
}

export default Layout