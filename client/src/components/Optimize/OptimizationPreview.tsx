import React, { useState } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import { RootState } from '../../store'
import { clearCurrentSession, applyFeedback, addFeedback, removeFeedback } from '../../store/slices/optimizationSlice'
import { AppDispatch } from '../../store'

const OptimizationPreview: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { currentSession, feedback, isProcessingFeedback } = useSelector((state: RootState) => state.optimization)
  const [selectedText, setSelectedText] = useState('')
  const [feedbackComment, setFeedbackComment] = useState('')
  const [isAddingFeedback, setIsAddingFeedback] = useState(false)

  const handleTextSelection = () => {
    const selection = window.getSelection()
    if (selection && selection.toString().trim()) {
      setSelectedText(selection.toString().trim())
      setIsAddingFeedback(true)
    }
  }

  const handleAddFeedback = () => {
    if (selectedText && feedbackComment) {
      dispatch(addFeedback({
        id: Date.now().toString(),
        session_id: currentSession!.id,
        section_highlight: selectedText,
        user_comment: feedbackComment,
        is_processed: false,
        created_at: new Date().toISOString()
      }))
      
      setSelectedText('')
      setFeedbackComment('')
      setIsAddingFeedback(false)
    }
  }

  const handleApplyFeedback = () => {
    if (feedback.length === 0) return

    dispatch(applyFeedback({
      sessionId: currentSession!.id,
      feedback: feedback.map(f => ({
        section_highlight: f.section_highlight,
        user_comment: f.user_comment
      }))
    }))
  }

  const handleDownload = () => {
    if (currentSession?.optimized_content) {
      const blob = new Blob([currentSession.optimized_content], { type: 'text/plain' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `optimized_resume_${currentSession.id}.txt`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    }
  }

  const handleStartOver = () => {
    dispatch(clearCurrentSession())
  }

  if (!currentSession) return null

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">Optimization Preview</h1>
              <p className="text-gray-600">
                Review your optimized resume. Select text to add feedback, then click "Apply Feedback" to refine.
              </p>
            </div>
            <div className="flex space-x-3">
              <button
                onClick={handleStartOver}
                className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Start Over
              </button>
              <button
                onClick={handleDownload}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              >
                <svg className="-ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                Download
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Optimized Resume Content */}
        <div className="lg:col-span-2">
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Optimized Resume</h3>
              <div className="bg-gray-50 rounded-lg p-6 max-h-96 overflow-y-auto">
                {currentSession.status === 'processing' ? (
                  <div className="flex items-center justify-center py-12">
                    <div className="text-center">
                      <svg className="animate-spin -ml-1 mr-3 h-8 w-8 text-blue-600 mx-auto" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      <p className="text-gray-600 mt-2">AI is optimizing your resume...</p>
                    </div>
                  </div>
                ) : (
                  <div
                    className="prose prose-sm max-w-none cursor-text"
                    onMouseUp={handleTextSelection}
                    style={{ userSelect: 'text' }}
                  >
                    <pre className="whitespace-pre-wrap font-sans text-sm text-gray-800">
                      {currentSession.optimized_content || 'No optimized content available'}
                    </pre>
                  </div>
                )}
              </div>
              <p className="text-xs text-gray-500 mt-2">
                Select any text above to add feedback for further optimization
              </p>
            </div>
          </div>
        </div>

        {/* Feedback Panel */}
        <div className="space-y-6">
          {/* Add Feedback */}
          {isAddingFeedback && (
            <div className="bg-white shadow rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <h4 className="text-md font-medium text-gray-900 mb-3">Add Feedback</h4>
                <div className="space-y-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Selected Text
                    </label>
                    <div className="p-2 bg-yellow-50 border border-yellow-200 rounded text-sm">
                      "{selectedText}"
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Your Feedback
                    </label>
                    <textarea
                      value={feedbackComment}
                      onChange={(e) => setFeedbackComment(e.target.value)}
                      placeholder="Explain how this section should be improved..."
                      rows={3}
                      className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 text-sm"
                    />
                  </div>
                  <div className="flex space-x-2">
                    <button
                      onClick={handleAddFeedback}
                      disabled={!feedbackComment.trim()}
                      className="flex-1 inline-flex justify-center py-2 px-3 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                    >
                      Add
                    </button>
                    <button
                      onClick={() => {
                        setIsAddingFeedback(false)
                        setSelectedText('')
                        setFeedbackComment('')
                      }}
                      className="flex-1 inline-flex justify-center py-2 px-3 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Feedback List */}
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h4 className="text-md font-medium text-gray-900 mb-3">
                Feedback ({feedback.length})
              </h4>
              
              {feedback.length === 0 ? (
                <p className="text-sm text-gray-500 text-center py-4">
                  No feedback added yet. Select text from the resume to add feedback.
                </p>
              ) : (
                <div className="space-y-3 max-h-64 overflow-y-auto">
                  {feedback.map((item) => (
                    <div key={item.id} className="border border-gray-200 rounded-lg p-3">
                      <div className="text-xs text-gray-500 mb-1">Selected text:</div>
                      <div className="text-sm bg-yellow-50 p-2 rounded mb-2">
                        "{item.section_highlight}"
                      </div>
                      <div className="text-xs text-gray-500 mb-1">Feedback:</div>
                      <div className="text-sm text-gray-800 mb-2">
                        {item.user_comment}
                      </div>
                      <button
                        onClick={() => dispatch(removeFeedback(item.id))}
                        className="text-xs text-red-600 hover:text-red-800"
                      >
                        Remove
                      </button>
                    </div>
                  ))}
                </div>
              )}

              {feedback.length > 0 && (
                <button
                  onClick={handleApplyFeedback}
                  disabled={isProcessingFeedback}
                  className="w-full mt-4 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                >
                  {isProcessingFeedback ? (
                    <>
                      <svg className="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Applying...
                    </>
                  ) : (
                    'Apply Feedback'
                  )}
                </button>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default OptimizationPreview