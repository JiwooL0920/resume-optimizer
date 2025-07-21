import React, { useState, useEffect } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState } from '../../store'
import { fetchProfile } from '../../store/slices/authSlice'
import { AppDispatch } from '../../store'
import Header from './Header'
import Login from '../Auth/Login'
import Dashboard from '../Dashboard/Dashboard'
import OptimizeResume from '../Optimize/OptimizeResume'
import Settings from '../Settings/Settings'

const Layout: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { isAuthenticated, token } = useSelector((state: RootState) => state.auth)
  const [currentPage, setCurrentPage] = useState('dashboard')

  useEffect(() => {
    if (token && !isAuthenticated) {
      dispatch(fetchProfile())
    }
  }, [dispatch, token, isAuthenticated])

  const handleNavigate = (page: string) => {
    setCurrentPage(page)
  }

  const renderCurrentPage = () => {
    // Temporarily disable auth check
    // if (!isAuthenticated) {
    //   return <Login />
    // }

    switch (currentPage) {
      case 'dashboard':
        return <Dashboard />
      case 'optimize':
        return <OptimizeResume />
      case 'settings':
        return <Settings />
      default:
        return <Dashboard />
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header onNavigate={handleNavigate} currentPage={currentPage} />
      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        {renderCurrentPage()}
      </main>
    </div>
  )
}

export default Layout