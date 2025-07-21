import { configureStore } from '@reduxjs/toolkit'
import authSlice from './slices/authSlice'
import resumeSlice from './slices/resumeSlice'
import optimizationSlice from './slices/optimizationSlice'
import apiKeysSlice from './slices/apiKeysSlice'

export const store = configureStore({
  reducer: {
    auth: authSlice,
    resume: resumeSlice,
    optimization: optimizationSlice,
    apiKeys: apiKeysSlice,
  },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch