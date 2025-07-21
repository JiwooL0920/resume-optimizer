import React, { useState, useEffect } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState } from '../../store'
import { optimizeResume, setSelectedAiModel, setKeepOnePage } from '../../store/slices/optimizationSlice'
import { AppDispatch } from '../../store'
import OptimizationPreview from './OptimizationPreview'

interface ApiKey {
  id: string
  provider: string
  masked_key: string
  created_at: string
}

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

  const [jobDescriptionUrl, setJobDescriptionUrl] = useState('')
  const [jobDescriptionText, setJobDescriptionText] = useState('')
  const [inputMode, setInputMode] = useState<'url' | 'text'>('url')
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([])
  const [selectedApiKeyId, setSelectedApiKeyId] = useState('')
  const [isLoadingApiKeys, setIsLoadingApiKeys] = useState(true)
  
  // Fetch user's API keys on component mount
  useEffect(() => {
    const fetchApiKeys = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/v1/user/api-keys', {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        if (response.ok) {
          const data = await response.json()
          setApiKeys(data.api_keys || [])
          // Auto-select first key if available and matches current AI model
          if (data.api_keys && data.api_keys.length > 0) {
            const matchingKey = data.api_keys.find((key: ApiKey) => {
              if (selectedAiModel.startsWith('gpt-') && key.provider === 'openai') return true
              if (selectedAiModel.startsWith('claude-') && key.provider === 'anthropic') return true
              return false
            })
            if (matchingKey) {
              setSelectedApiKeyId(matchingKey.id)
            }
          }
        }
      } catch (error) {
        console.error('Failed to fetch API keys:', error)
      } finally {
        setIsLoadingApiKeys(false)
      }
    }
    
    fetchApiKeys()
  }, [selectedAiModel])
  
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
      alert('Please select a resume first from the Dashboard')
      return
    }

    if (!jobDescriptionUrl && !jobDescriptionText) {
      alert('Please provide either a job description URL or text')
      return
    }

    if (!selectedApiKeyId) {
      alert('Please select an API key or add one in Settings')
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
                <option value="gpt-4">GPT-4 (OpenAI)</option>
                <option value="gpt-3.5-turbo">GPT-3.5 Turbo (OpenAI)</option>
                <option value="claude-3-opus">Claude 3 Opus (Anthropic)</option>
                <option value="claude-3-sonnet">Claude 3 Sonnet (Anthropic)</option>
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
                  <a 
                    href="/settings" 
                    className="text-sm text-blue-600 hover:text-blue-500"
                  >
                    Add API keys in Settings →
                  </a>
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
            <div className="mt-4 rounded-md bg-red-50 p-4">
              <div className="flex">
                <div className="flex-shrink-0">
                  <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                  </svg>
                </div>
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-red-800">Optimization failed</h3>
                  <div className="mt-2 text-sm text-red-700">
                    <p>{error}</p>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default OptimizeResume