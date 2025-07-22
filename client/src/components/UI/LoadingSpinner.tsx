import React from 'react'

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg'
  color?: 'blue' | 'gray' | 'white'
  text?: string
  className?: string
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  color = 'blue',
  text,
  className = ''
}) => {
  const getSizeClasses = () => {
    switch (size) {
      case 'sm':
        return 'h-4 w-4'
      case 'lg':
        return 'h-8 w-8'
      case 'md':
      default:
        return 'h-5 w-5'
    }
  }

  const getColorClasses = () => {
    switch (color) {
      case 'gray':
        return 'text-gray-600'
      case 'white':
        return 'text-white'
      case 'blue':
      default:
        return 'text-blue-600'
    }
  }

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <svg
        className={`animate-spin ${getSizeClasses()} ${getColorClasses()}`}
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          className="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          strokeWidth="4"
        />
        <path
          className="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        />
      </svg>
      {text && (
        <span className={`ml-2 text-sm ${getColorClasses()}`}>
          {text}
        </span>
      )}
    </div>
  )
}

export default LoadingSpinner