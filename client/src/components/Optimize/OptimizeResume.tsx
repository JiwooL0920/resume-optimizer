import React, { useState, useEffect } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { Link } from 'react-router-dom'
import { RootState } from '../../store'
import { optimizeResume, setSelectedAiModel, setKeepOnePage } from '../../store/slices/optimizationSlice'
import { fetchApiKeys } from '../../store/slices/apiKeysSlice'
import { AppDispatch } from '../../store'
import { ApiKey } from '../../types'
import { AI_MODELS } from '../../config/api'
import { ErrorMessage } from '../UI'
import OptimizationPreview from './OptimizationPreview'

const OptimizeResume: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { selectedResume } = useSelector((state: RootState) => state.resume)
  const { 
    currentSession, 
    isOptimizing, 
    error, 
    selectedAiModel, 
    keepOnePage 
  } = useSelector((state: RootState) => state.optimization)
  const { keys: apiKeys, isLoading: isLoadingApiKeys } = useSelector((state: RootState) => state.apiKeys)

  const [jobDescriptionUrl, setJobDescriptionUrl] = useState('')
  const [jobDescriptionText, setJobDescriptionText] = useState('')
  const [inputMode, setInputMode] = useState<'url' | 'text'>('url')
  const [selectedApiKeyId, setSelectedApiKeyId] = useState('')
  
  // Fetch user's API keys on component mount
  useEffect(() => {
    dispatch(fetchApiKeys())
  }, [dispatch])
  
  // Update selected API key when AI model changes
  useEffect(() => {
    const matchingKey = apiKeys.find(key => {
      if (selectedAiModel.startsWith('gpt-') && key.provider === 'openai') return true
      if (selectedAiModel.startsWith('claude-') && key.provider === 'anthropic') return true
      return false
    })
    if (matchingKey) {
      setSelectedApiKeyId(matchingKey.id)
    } else {
      setSelectedApiKeyId('')
    }
  }, [selectedAiModel, apiKeys])

  const handleOptimize = () => {
    if (!selectedResume) {
      // This is handled by the disabled state, but keeping as fallback
      return
    }

    if (!jobDescriptionUrl && !jobDescriptionText) {
      // This is handled by the disabled state, but keeping as fallback
      return
    }

    if (!selectedApiKeyId) {
      // This is handled by the disabled state, but keeping as fallback
      return
    }

    dispatch(optimizeResume({
      resumeId: selectedResume.id,
      jobDescriptionUrl: inputMode === 'url' ? jobDescriptionUrl : undefined,
      jobDescriptionText: inputMode === 'text' ? jobDescriptionText : undefined,
      aiModel: selectedAiModel,
      keepOnePage,
      userApiKey: selectedApiKeyId
    }))
  }

  if (currentSession) {
    return <OptimizationPreview />
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Optimize Resume</h1>
          <p className="text-gray-600">
            Upload a resume from the Dashboard, then provide a job description to optimize your resume using AI
          </p>
        </div>
      </div>

      {/* Selected Resume */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Selected Resume</h3>
          {selectedResume ? (
            <div className="flex items-center p-4 border border-green-200 bg-green-50 rounded-lg">
              <div className="flex-shrink-0">
                <svg className="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div className="ml-3">
                <h4 className="text-sm font-medium text-green-800">{selectedResume.title}</h4>
                <p className="text-sm text-green-700">
                  {selectedResume.file_type.toUpperCase()} • Uploaded {new Date(selectedResume.created_at).toLocaleDateString()}
                </p>
              </div>
            </div>
          ) : (
            <div className="text-center py-6 border-2 border-dashed border-gray-300 rounded-lg">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p className="mt-2 text-sm text-gray-500">No resume selected</p>
              <p className="text-xs text-gray-400">Go to Dashboard to select a resume</p>
            </div>
          )}
        </div>
      </div>

      {/* Job Description */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Job Description</h3>
          
          <div className="mb-4">
            <div className="sm:hidden">
              <select
                value={inputMode}
                onChange={(e) => setInputMode(e.target.value as 'url' | 'text')}
                className="block w-full rounded-md border-gray-300 focus:border-blue-500 focus:ring-blue-500"
              >
                <option value="url">URL</option>
                <option value="text">Text</option>
              </select>
            </div>
            <div className="hidden sm:block">
              <nav className="flex space-x-8">
                <button
                  onClick={() => setInputMode('url')}
                  className={`${
                    inputMode === 'url'
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm`}
                >
                  Job URL
                </button>
                <button
                  onClick={() => setInputMode('text')}
                  className={`${
                    inputMode === 'text'
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm`}
                >
                  Paste Text
                </button>
              </nav>
            </div>
          </div>

          {inputMode === 'url' ? (
            <div>
              <label htmlFor="job-url" className="block text-sm font-medium text-gray-700 mb-2">
                Job Description URL
              </label>
              <input
                type="url"
                id="job-url"
                value={jobDescriptionUrl}
                onChange={(e) => setJobDescriptionUrl(e.target.value)}
                placeholder="https://example.com/job-posting"
                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              />
              <p className="mt-1 text-xs text-gray-500">
                We'll extract the job description from this URL
              </p>
            </div>
          ) : (
            <div>
              <label htmlFor="job-text" className="block text-sm font-medium text-gray-700 mb-2">
                Job Description Text
              </label>
              <textarea
                id="job-text"
                rows={8}
                value={jobDescriptionText}
                onChange={(e) => setJobDescriptionText(e.target.value)}
                placeholder="Paste the job description here..."
                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              />
            </div>
          )}
        </div>
      </div>

      {/* AI Settings */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">AI Settings</h3>
          
          <div className="space-y-4">
            <div>
              <label htmlFor="ai-model" className="block text-sm font-medium text-gray-700 mb-2">
                AI Model
              </label>
              <select
                id="ai-model"
                value={selectedAiModel}
                onChange={(e) => dispatch(setSelectedAiModel(e.target.value))}
                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              >
                {Object.entries(AI_MODELS).map(([key, model]) => (
                  <option key={key} value={key}>
                    {model.label} ({model.provider.charAt(0).toUpperCase() + model.provider.slice(1)})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="api-key" className="block text-sm font-medium text-gray-700 mb-2">
                API Key
              </label>
              {isLoadingApiKeys ? (
                <div className="animate-pulse h-10 bg-gray-200 rounded-md"></div>
              ) : apiKeys.length === 0 ? (
                <div className="text-center py-4 border-2 border-dashed border-gray-300 rounded-lg">
                  <svg className="mx-auto h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                  </svg>
                  <p className="mt-2 text-sm text-gray-500">No API keys found</p>
                  <Link 
                    to="/settings" 
                    className="text-sm text-blue-600 hover:text-blue-500"
                  >
                    Add API keys in Settings →
                  </Link>
                </div>
              ) : (
                <select
                  id="api-key"
                  value={selectedApiKeyId}
                  onChange={(e) => setSelectedApiKeyId(e.target.value)}
                  className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                >
                  <option value="">Select an API key...</option>
                  {apiKeys
                    .filter(key => {
                      if (selectedAiModel.startsWith('gpt-')) return key.provider === 'openai'
                      if (selectedAiModel.startsWith('claude-')) return key.provider === 'anthropic'
                      return true
                    })
                    .map(key => (
                      <option key={key.id} value={key.id}>
                        {key.provider.charAt(0).toUpperCase() + key.provider.slice(1)} - {key.masked_key}
                      </option>
                    ))
                  }
                </select>
              )}
              {apiKeys.length > 0 && !selectedApiKeyId && (
                <p className="text-xs text-red-500 mt-1">
                  Please select an API key to continue
                </p>
              )}
            </div>

            <div className="flex items-start">
              <div className="flex items-center h-5">
                <input
                  id="keep-one-page"
                  type="checkbox"
                  checked={keepOnePage}
                  onChange={(e) => dispatch(setKeepOnePage(e.target.checked))}
                  className="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
                />
              </div>
              <div className="ml-3 text-sm">
                <label htmlFor="keep-one-page" className="font-medium text-gray-700">
                  Keep resume within one page
                </label>
                <p className="text-gray-500">
                  Optimize for single-page format (recommended for most positions)
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Action Button */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <button
            onClick={handleOptimize}
            disabled={isOptimizing || !selectedResume || (!jobDescriptionUrl && !jobDescriptionText) || !selectedApiKeyId}
            className="w-full flex justify-center py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isOptimizing ? (
              <>
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Optimizing Resume...
              </>
            ) : (
              'Optimize Resume'
            )}
          </button>
          
          {error && (
            <ErrorMessage 
              error={error} 
              title="Optimization failed"
              className="mt-4"
            />
          )}
        </div>
      </div>
    </div>
  )
}

export default OptimizeResume