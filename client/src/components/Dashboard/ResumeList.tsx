import React, { memo, useCallback } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState } from '../../store'
import { selectResume, deleteResume } from '../../store/slices/resumeSlice'
import { AppDispatch } from '../../store'
import { Resume } from '../../types'

const ResumeList: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { resumes, selectedResume } = useSelector((state: RootState) => state.resume)

  const handleSelectResume = useCallback((resume: Resume) => {
    dispatch(selectResume(resume))
  }, [dispatch])

  const handleDeleteResume = useCallback((resumeId: string) => {
    if (window.confirm('Are you sure you want to delete this resume?')) {
      dispatch(deleteResume(resumeId))
    }
  }, [dispatch])

  const formatFileSize = useCallback((bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }, [])

  const formatDate = useCallback((dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }, [])

  if (resumes.length === 0) {
    return (
      <div className="text-center py-8">
        <svg
          className="mx-auto h-12 w-12 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">No resumes</h3>
        <p className="mt-1 text-sm text-gray-500">Get started by uploading your first resume.</p>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {resumes.map((resume) => (
        <div
          key={resume.id}
          className={`border rounded-lg p-4 cursor-pointer transition-all ${
            selectedResume?.id === resume.id
              ? 'border-blue-500 bg-blue-50 ring-2 ring-blue-200'
              : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
          }`}
          onClick={() => handleSelectResume(resume)}
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="flex-shrink-0">
                <div className="w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center">
                  <svg className="w-6 h-6 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-center space-x-2">
                  <h4 className="text-sm font-medium text-gray-900 truncate">
                    {resume.title}
                  </h4>
                  {resume.is_active && (
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                      Active
                    </span>
                  )}
                  {selectedResume?.id === resume.id && (
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                      Selected
                    </span>
                  )}
                </div>
                <div className="flex items-center space-x-4 mt-1">
                  <p className="text-sm text-gray-500 uppercase">
                    {resume.file_type}
                  </p>
                  {resume.file_size && (
                    <p className="text-sm text-gray-500">
                      {formatFileSize(resume.file_size)}
                    </p>
                  )}
                  <p className="text-sm text-gray-500">
                    Uploaded {formatDate(resume.created_at)}
                  </p>
                </div>
              </div>
            </div>
            
            <div className="flex items-center space-x-2">
              <button
                onClick={(e) => {
                  e.stopPropagation()
                  handleDeleteResume(resume.id)
                }}
                className="p-2 text-gray-400 hover:text-red-600 transition-colors"
                title="Delete resume"
              >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
              
              {selectedResume?.id === resume.id && (
                <div className="flex items-center">
                  <svg className="w-4 h-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

export default memo(ResumeList)