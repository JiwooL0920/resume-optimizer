import React from 'react'
import { Provider } from 'react-redux'
import { BrowserRouter } from 'react-router-dom'
import { store } from './store'
import { ErrorBoundary } from './components/ErrorBoundary'
import Layout from './components/Layout/Layout'

function App() {
  return (
    <ErrorBoundary>
      <Provider store={store}>
        <BrowserRouter>
          <Layout />
        </BrowserRouter>
      </Provider>
    </ErrorBoundary>
  )
}

export default App
