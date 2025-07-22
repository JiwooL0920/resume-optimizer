import React, { useState, useEffect } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState, AppDispatch } from '../../store'
import { fetchApiKeys, createApiKey, deleteApiKey, clearError } from '../../store/slices/apiKeysSlice'
import { ApiKey } from '../../types'
import { ErrorMessage, SuccessMessage } from '../UI'

const Settings: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { user } = useSelector((state: RootState) => state.auth)
  const { keys: apiKeys, isLoading, error } = useSelector((state: RootState) => state.apiKeys)
  const [isAddingKey, setIsAddingKey] = useState(false)
  const [newKeyProvider, setNewKeyProvider] = useState('openai')
  const [newKeyValue, setNewKeyValue] = useState('')
  const [showError, setShowError] = useState(false)
  const [showSuccess, setShowSuccess] = useState(false)
  const [successMessage, setSuccessMessage] = useState('')

  useEffect(() => {
    dispatch(fetchApiKeys())
  }, [dispatch])

  useEffect(() => {
    if (error) {
      setShowError(true)
    }
  }, [error])

  const handleAddApiKey = async () => {
    if (!newKeyValue.trim()) return
    const result = await dispatch(createApiKey({ provider: newKeyProvider, api_key: newKeyValue }))
    if (createApiKey.fulfilled.match(result)) {
      setSuccessMessage(`${newKeyProvider.charAt(0).toUpperCase() + newKeyProvider.slice(1)} API key added successfully`)
      setShowSuccess(true)
      setNewKeyValue('')
      setIsAddingKey(false)
    }
  }

  const handleDeleteApiKey = async (keyId: string, provider: string) => {
    if (!window.confirm('Are you sure you want to delete this API key?')) return
    const result = await dispatch(deleteApiKey(keyId))
    if (deleteApiKey.fulfilled.match(result)) {
      setSuccessMessage(`${provider.charAt(0).toUpperCase() + provider.slice(1)} API key deleted successfully`)
      setShowSuccess(true)
    }
  }

  const providerLogos: { [key: string]: string } = {
    openai: '🤖',
    anthropic: '🔮',
    google: '🧠',
    cohere: '🚀'
  }

  const providerNames: { [key: string]: string } = {
    openai: 'OpenAI',
    anthropic: 'Anthropic',
    google: 'Google AI',
    cohere: 'Cohere'
  }

  return (
    <div className="space-y-6">
      {/* Error Messages */}
      {showError && error && (
        <ErrorMessage 
          error={error} 
          onDismiss={() => {
            setShowError(false)
            dispatch(clearError())
          }}
        />
      )}

      {/* Success Messages */}
      {showSuccess && successMessage && (
        <SuccessMessage 
          message={successMessage} 
          onDismiss={() => {
            setShowSuccess(false)
            setSuccessMessage('')
          }}
        />
      )}

      {/* Header */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Settings</h1>
          <p className="text-gray-600">
            Manage your AI API keys and account preferences
          </p>
        </div>
      </div>

      {/* Profile Information */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Profile Information</h3>
          <div className="flex items-center space-x-4">
            {user?.picture_url && (
              <img
                className="h-16 w-16 rounded-full"
                src={user.picture_url}
                alt={user.name}
              />
            )}
            <div>
              <h4 className="text-lg font-medium text-gray-900">{user?.name}</h4>
              <p className="text-gray-600">{user?.email}</p>
              <p className="text-sm text-gray-500">
                Account created via Google OAuth
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* API Keys Management */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex justify-between items-center mb-4">
            <div>
              <h3 className="text-lg font-medium text-gray-900">AI API Keys</h3>
              <p className="text-sm text-gray-600">
                Add your AI service API keys to enable resume optimization
              </p>
            </div>
            {!isAddingKey && (
              <button
                onClick={() => setIsAddingKey(true)}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <svg className="-ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                Add API Key
              </button>
            )}
          </div>

          {/* Add API Key Form */}
          {isAddingKey && (
            <div className="mb-6 p-4 border border-gray-200 rounded-lg bg-gray-50">
              <h4 className="text-md font-medium text-gray-900 mb-3">Add New API Key</h4>
              <div className="space-y-4">
                <div>
                  <label htmlFor="provider" className="block text-sm font-medium text-gray-700 mb-2">
                    AI Provider
                  </label>
                  <select
                    id="provider"
                    value={newKeyProvider}
                    onChange={(e) => setNewKeyProvider(e.target.value)}
                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  >
                    <option value="openai">OpenAI (GPT-4, GPT-3.5)</option>
                    <option value="anthropic">Anthropic (Claude)</option>
                    <option value="google">Google AI (Gemini)</option>
                    <option value="cohere">Cohere</option>
                  </select>
                </div>
                <div>
                  <label htmlFor="api-key" className="block text-sm font-medium text-gray-700 mb-2">
                    API Key
                  </label>
                  <input
                    type="password"
                    id="api-key"
                    value={newKeyValue}
                    onChange={(e) => setNewKeyValue(e.target.value)}
                    placeholder="sk-... or your API key"
                    className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    Your API key will be encrypted and stored securely
                  </p>
                </div>
                <div className="flex space-x-3">
                  <button
                    onClick={handleAddApiKey}
                    disabled={isLoading || !newKeyValue.trim()}
                    className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                  >
                    {isLoading ? (
                      <>
                        <svg className="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Adding...
                      </>
                    ) : (
                      'Add Key'
                    )}
                  </button>
                  <button
                    onClick={() => {
                      setIsAddingKey(false)
                      setNewKeyValue('')
                    }}
                    className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* API Keys List */}
          {apiKeys.length === 0 ? (
            <div className="text-center py-8">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
              <h3 className="mt-2 text-sm font-medium text-gray-900">No API keys</h3>
              <p className="mt-1 text-sm text-gray-500">
                Add your AI service API keys to start optimizing resumes.
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {apiKeys.map((key) => (
                <div key={key.id} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                  <div className="flex items-center space-x-3">
                    <div className="text-2xl">
                      {providerLogos[key.provider] || '🔑'}
                    </div>
                    <div>
                      <h4 className="text-sm font-medium text-gray-900">
                        {providerNames[key.provider] || key.provider}
                      </h4>
                      <p className="text-sm text-gray-500">
                        {key.masked_key} • Added {new Date(key.created_at).toLocaleDateString()}
                      </p>
                    </div>
                  </div>
                  <button
                    onClick={() => handleDeleteApiKey(key.id, key.provider)}
                    className="text-red-600 hover:text-red-800 p-2"
                    title="Delete API key"
                  >
                    <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Instructions */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-blue-900 mb-3">How to get API keys</h3>
          <div className="space-y-3 text-sm text-blue-800">
            <div>
              <strong>OpenAI:</strong> Visit{' '}
              <a href="https://platform.openai.com/api-keys" target="_blank" rel="noopener noreferrer" className="underline">
                platform.openai.com/api-keys
              </a>{' '}
              to create your API key
            </div>
            <div>
              <strong>Anthropic:</strong> Visit{' '}
              <a href="https://console.anthropic.com/" target="_blank" rel="noopener noreferrer" className="underline">
                console.anthropic.com
              </a>{' '}
              to get your Claude API key
            </div>
            <div>
              <strong>Google AI:</strong> Visit{' '}
              <a href="https://ai.google.dev/" target="_blank" rel="noopener noreferrer" className="underline">
                ai.google.dev
              </a>{' '}
              to get started with Gemini API
            </div>
          </div>
          <p className="text-xs text-blue-700 mt-4">
            <strong>Security:</strong> Your API keys are encrypted before storage and never shared. 
            You can delete them at any time.
          </p>
        </div>
      </div>
    </div>
  )
}

export default Settings