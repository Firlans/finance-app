import { Config } from './Config.js'
import { Notification } from './Notification.js'

async function getCurrentUser(token) {
  try {
    const response = await fetch(`${Config.url}/users/current`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Failed to fetch user data: ${response.status} ${response.statusText} ${errorText}`)
    }

    const json = await response.json()
    return json?.data
  } catch (err) {
    console.error('Error fetching current user:', err)
    Notification.show('Error fetching current user')
    return null
  }
}
async function getAccounts() {
}
async function getTransactions() {
}

export { getCurrentUser, getAccounts, getTransactions }
