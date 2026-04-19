export const parseApiError = async (response) => {
  const data = await response.json().catch(() => null)
  const message = data?.error || data?.message || 'Request failed'
  const parts = []

  if (typeof data?.retry_after_minutes === 'number') {
    parts.push(`${data.retry_after_minutes} minute${data.retry_after_minutes === 1 ? '' : 's'}`)
  }
  if (typeof data?.retry_after_seconds === 'number') {
    parts.push(`${data.retry_after_seconds} second${data.retry_after_seconds === 1 ? '' : 's'}`)
  }

  if (parts.length) {
    return `${message} Please try again after ${parts.join(' / ')}.`
  }

  return message
}
