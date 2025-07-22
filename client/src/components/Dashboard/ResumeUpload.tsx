import React, { useState, useRef, memo, useCallback } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState } from '../../store'
import { uploadResume } from '../../store/slices/resumeSlice'
import { AppDispatch } from '../../store'
import { ErrorMessage } from '../UI'

const ResumeUpload: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { isLoading, uploadProgress, error } = useSelector((state: RootState) => state.resume)
  const [dragActive, setDragActive] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  const handleFiles = useCallback((files: FileList | null) => {
    if (files && files[0]) {
      const file = files[0]
      
      // Validation is now handled by the file input accept attribute and backend
      // Additional client-side validation can be added here if needed

      dispatch(uploadResume(file))
    }
  }, [dispatch])

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true)
    } else if (e.type === 'dragleave') {
      setDragActive(false)
    }
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFiles(e.dataTransfer.files)
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.preventDefault()
    handleFiles(e.target.files)
  }

  const handleButtonClick = () => {
    inputRef.current?.click()
  }

  return (
    <div className="space-y-4">
      <div
        className={`relative border-2 border-dashed rounded-lg p-6 ${
          dragActive
            ? 'border-blue-400 bg-blue-50'
            : 'border-gray-300 hover:border-gray-400'
        } ${isLoading ? 'opacity-50 pointer-events-none' : ''}`}
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
      >
        <div className="text-center">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            stroke="currentColor"
            fill="none"
            viewBox="0 0 48 48"
          >
            <path
              d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
              strokeWidth={2}
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
          <div className="mt-4">
            <label htmlFor="file-upload" className="cursor-pointer">
              <span className="mt-2 block text-sm font-medium text-gray-900">
                {isLoading ? 'Uploading...' : 'Drop your resume here or click to browse'}
              </span>
              <span className="mt-1 block text-xs text-gray-500">
                PDF, DOC, DOCX up to 10MB
              </span>
            </label>
            <input
              ref={inputRef}
              id="file-upload"
              name="file-upload"
              type="file"
              className="sr-only"
              accept=".pdf,.doc,.docx"
              onChange={handleChange}
              disabled={isLoading}
            />
          </div>
          
          {!isLoading && (
            <button
              type="button"
              onClick={handleButtonClick}
              className="mt-4 inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <svg className="-ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              Select File
            </button>
          )}
        </div>

        {isLoading && uploadProgress > 0 && (
          <div className="mt-4">
            <div className="bg-gray-200 rounded-full h-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              ></div>
            </div>
            <p className="text-sm text-gray-600 mt-1">{uploadProgress}% uploaded</p>
          </div>
        )}
      </div>

      {error && (
        <ErrorMessage 
          error={error} 
          title="Upload failed"
        />
      )}
    </div>
  )
}

export default memo(ResumeUpload)