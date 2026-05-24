import { Config } from './Config.js'
import { parseApiError } from '@packages/utils/Error.js'

const API_BASE = import.meta.env.VITE_BACKEND_SERVICE || Config.url

const normalizeToken = (token) => {
  if (!token) return ''
  return String(token)
    .trim()
    .replace(/^Bearer\s+/i, '')
}

const getAuthHeaders = (token, json = true) => {
  const headers = {}
  const normalizedToken = normalizeToken(token)
  if (normalizedToken) {
    headers.Authorization = `Bearer ${normalizedToken}`
  }
  if (json) {
    headers['Content-Type'] = 'application/json'
  }
  return headers
}

const request = async (path, options = {}) => {
  const response = await fetch(`${API_BASE}${path}`, options)

  if (!response.ok) {
    const errorMessage = await parseApiError(response)
    throw new Error(errorMessage)
  }

  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    return await response.json()
  }

  return null
}

async function getCurrentUser(token) {
  if (!token) {
    return null
  }

  const json = await request('/users/current', {
    headers: getAuthHeaders(token, false)
  })
  return json?.data ?? null
}

async function getAccounts(token) {
  const json = await request('/accounts', {
    headers: getAuthHeaders(token, false)
  })
  return json?.data || []
}

async function createAccount(token, payload) {
  const json = await request('/accounts', {
    method: 'POST',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function updateAccount(token, id, payload) {
  const json = await request(`/accounts/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function deleteAccount(token, id) {
  const json = await request(`/accounts/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(token, false)
  })
  return json?.data || null
}

async function getTransactions(token) {
  const json = await request('/transactions', {
    headers: getAuthHeaders(token, false)
  })
  return json?.data || []
}

async function createTransaction(token, payload) {
  const json = await request('/transactions', {
    method: 'POST',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function updateTransaction(token, id, payload) {
  const json = await request(`/transactions/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function deleteTransaction(token, id) {
  const json = await request(`/transactions/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(token, false)
  })
  return json?.data || null
}

async function logout(token) {
  const json = await request('/users/logout', {
    method: 'POST',
    headers: getAuthHeaders(token, false)
  })
  return json || null
}

async function getCategories(token) {
  const json = await request('/categories', {
    headers: getAuthHeaders(token, false)
  })
  return json?.data || []
}

async function createCategory(token, payload) {
  const json = await request('/categories', {
    method: 'POST',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function updateCategory(token, id, payload) {
  const json = await request(`/categories/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function deleteCategory(token, id) {
  const json = await request(`/categories/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(token, false)
  })
  return json?.data || null
}

export {
  getCurrentUser,
  getAccounts,
  createAccount,
  updateAccount,
  deleteAccount,
  getTransactions,
  createTransaction,
  updateTransaction,
  deleteTransaction,
  logout,
  getCategories,
  createCategory,
  updateCategory,
  deleteCategory
}
