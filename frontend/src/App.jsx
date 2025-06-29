import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { AuthProvider } from './contexts/AuthContext'
import { ErrorBoundary } from './components'
import Layout from './components/Layout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import Categories from './pages/Categories'
import Inbound from './pages/Inbound'
import Outbound from './pages/Outbound'
import Inventory from './pages/Inventory'
import Reports from './pages/Reports'
import Users from './pages/Users'
import PrivateRoute from './components/PrivateRoute'
import { handleErrorBoundaryError } from './utils/errorHandler'
import './App.css'

function App() {
  return (
    <ErrorBoundary 
      onRetry={() => window.location.reload()}
      showErrorDetails={process.env.NODE_ENV === 'development'}
    >
      <ConfigProvider 
        locale={zhCN}
        theme={{
          token: {
            colorPrimary: '#1890ff',
            borderRadius: 8,
          },
        }}
      >
        <AuthProvider>
          <Router>
            <div className="App">
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/" element={
                  <PrivateRoute>
                    <ErrorBoundary>
                      <Layout />
                    </ErrorBoundary>
                  </PrivateRoute>
                }>
                  <Route index element={<Navigate to="/dashboard" replace />} />
                  <Route path="dashboard" element={
                    <ErrorBoundary>
                      <Dashboard />
                    </ErrorBoundary>
                  } />
                  <Route path="categories" element={
                    <ErrorBoundary>
                      <Categories />
                    </ErrorBoundary>
                  } />
                  <Route path="inbound" element={
                    <ErrorBoundary>
                      <Inbound />
                    </ErrorBoundary>
                  } />
                  <Route path="outbound" element={
                    <ErrorBoundary>
                      <Outbound />
                    </ErrorBoundary>
                  } />
                  <Route path="inventory" element={
                    <ErrorBoundary>
                      <Inventory />
                    </ErrorBoundary>
                  } />
                  <Route path="reports" element={
                    <ErrorBoundary>
                      <Reports />
                    </ErrorBoundary>
                  } />
                  <Route path="users" element={
                    <ErrorBoundary>
                      <Users />
                    </ErrorBoundary>
                  } />
                </Route>
              </Routes>
            </div>
          </Router>
        </AuthProvider>
      </ConfigProvider>
    </ErrorBoundary>
  )
}

export default App