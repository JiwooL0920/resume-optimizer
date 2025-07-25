import React, { memo, useCallback } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { Link, useNavigate } from 'react-router-dom'
import { RootState } from '../../store'
import { logout } from '../../store/slices/authSlice'
import { AppDispatch } from '../../store'

interface HeaderProps {
  currentPage: string
}

const Header: React.FC<HeaderProps> = ({ currentPage }) => {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const { user, isAuthenticated } = useSelector((state: RootState) => state.auth)

  const handleLogout = useCallback(() => {
    dispatch(logout())
    navigate('/login')
  }, [dispatch, navigate])

  return (
    <header className="bg-white shadow-sm border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <h1 className="text-2xl font-bold text-gray-900">ResumeOptimizer</h1>
            </div>
            
            {isAuthenticated && (
              <nav className="hidden md:ml-6 md:flex md:space-x-8">
                <Link
                  to="/dashboard"
                  className={`${
                    currentPage === 'dashboard'
                      ? 'border-blue-500 text-gray-900'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm`}
                >
                  Dashboard
                </Link>
                <Link
                  to="/optimize"
                  className={`${
                    currentPage === 'optimize'
                      ? 'border-blue-500 text-gray-900'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm`}
                >
                  Optimize
                </Link>
                <Link
                  to="/settings"
                  className={`${
                    currentPage === 'settings'
                      ? 'border-blue-500 text-gray-900'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm`}
                >
                  Settings
                </Link>
              </nav>
            )}
          </div>

          <div className="flex items-center">
            {isAuthenticated ? (
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2">
                  {user?.picture_url && (
                    <img
                      className="h-8 w-8 rounded-full"
                      src={user.picture_url}
                      alt={user.name}
                    />
                  )}
                  <span className="text-sm text-gray-700">{user?.name}</span>
                </div>
                <button
                  onClick={handleLogout}
                  className="bg-gray-100 hover:bg-gray-200 px-3 py-2 rounded-md text-sm font-medium text-gray-700"
                >
                  Sign out
                </button>
              </div>
            ) : (
              <Link
                to="/login"
                className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
              >
                Sign in
              </Link>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}

export default memo(Header)