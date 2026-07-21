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

async function refreshToken(token) {
  if (!token) {
    return null
  }

  const json = await request('/users/refresh-token', {
    method: 'POST',
    headers: getAuthHeaders(token, false)
  })

  return json?.data?.access_token || json?.access_token || null
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

async function getTransactions(token, from = null, to = null, page = null) {
  let url = '/transactions'
  const params = new URLSearchParams()

  // Tambahkan parameter jika nilai from/to dikirimkan
  if (from) params.append('from', from)
  if (to) params.append('to', to)
  if (page) params.append('page', page)

  // Gabungkan URL dengan query string jika ada
  const queryString = params.toString()
  if (queryString) {
    url += `?${queryString}`
  }

  const json = await request(url, {
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

async function getSummary(token, module, from = null, to = null) {
  let url = '/summary'
  const params = new URLSearchParams()
  if (module) params.append('module', module)
  if (from) params.append('from', from)
  if (to) params.append('to', to)

  const queryString = params.toString()
  if (queryString) {
    url += `?${queryString}`
  }

  const json = await request(url, {
    headers: getAuthHeaders(token, false)
  })

  return json?.data || []
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

// ------------------------
// Loans
// ------------------------
async function getLoans(token) {
  const json = await request('/loans', {
    headers: getAuthHeaders(token, false)
  })
  return json?.data || []
}

async function createLoan(token, payload) {
  const json = await request('/loans', {
    method: 'POST',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function updateLoan(token, id, payload) {
  const json = await request(`/loans/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function deleteLoan(token, id) {
  const json = await request(`/loans/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(token, false)
  })
  return json?.data || null
}

// ------------------------
// Payments
// ------------------------
async function getPaymentsByLoan(token, loanId) {
  const url = `/payments?loan_id=${encodeURIComponent(loanId)}`
  const json = await request(url, {
    headers: getAuthHeaders(token, false)
  })
  return json?.data || []
}

async function createPayment(token, payload) {
  const json = await request('/payments', {
    method: 'POST',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function updatePayment(token, id, payload) {
  const json = await request(`/payments/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function deletePayment(token, id) {
  const json = await request(`/payments/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(token, false)
  })
  return json?.data || null
}

// ------------------------
// Budgets
// ------------------------
async function getBudgets(token, date = null) {
  let url = '/budgets/summary'
  if (date) {
    url += `?date=${date}`
  }
  const json = await request(url, {
    headers: getAuthHeaders(token, false)
  })
  return json?.data || []
}

async function upsertBudget(token, payload) {
  const json = await request('/budgets', {
    method: 'POST',
    headers: getAuthHeaders(token),
    body: JSON.stringify(payload)
  })
  return json?.data || null
}

async function deleteBudget(token, id) {
  const json = await request(`/budgets/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(token, false)
  })
  return json?.data || null
}

export {
  getCurrentUser,
  refreshToken,
  getAccounts,
  createAccount,
  updateAccount,
  deleteAccount,
  getTransactions,
  createTransaction,
  updateTransaction,
  deleteTransaction,
  getSummary,
  logout,
  getCategories,
  createCategory,
  updateCategory,
  deleteCategory,

  // loans
  getLoans,
  createLoan,
  updateLoan,
  deleteLoan,

  // payments
  getPaymentsByLoan,
  createPayment,
  updatePayment,
  deletePayment,

  // budgets
  getBudgets,
  upsertBudget,
  deleteBudget
}

