const BASE_URL = 'http://localhost:8080/api'

async function request(url, options = {}) {
  const token = localStorage.getItem('token')
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${url}`, { ...options, headers })
  const data = await res.json()

  if (!res.ok) {
    throw new Error(data.error || '请求失败')
  }
  return data
}

export const authApi = {
  register: (data) => request('/auth/register', { method: 'POST', body: JSON.stringify(data) }),
  login: (data) => request('/auth/login', { method: 'POST', body: JSON.stringify(data) }),
  getMe: () => request('/auth/me'),
}

export const categoryApi = {
  getAll: () => request('/categories'),
  create: (data) => request('/admin/categories', { method: 'POST', body: JSON.stringify(data) }),
  update: (id, data) => request(`/admin/categories/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  delete: (id) => request(`/admin/categories/${id}`, { method: 'DELETE' }),
}

export const articleApi = {
  getList: (params) => {
    const query = new URLSearchParams(params).toString()
    return request(`/articles?${query}`)
  },
  getDetail: (id) => request(`/articles/${id}`),
  create: (data) => request('/articles', { method: 'POST', body: JSON.stringify(data) }),
  update: (id, data) => request(`/articles/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  delete: (id) => request(`/articles/${id}`, { method: 'DELETE' }),
  getMine: (params) => {
    const query = new URLSearchParams(params).toString()
    return request(`/articles/mine?${query}`)
  },
  search: (params) => {
    const query = new URLSearchParams(params).toString()
    return request(`/articles/search?${query}`)
  },
}

export const versionApi = {
  getList: (articleId, params) => {
    const query = params ? new URLSearchParams(params).toString() : ''
    return request(`/articles/${articleId}/versions${query ? '?' + query : ''}`)
  },
  getDetail: (versionId) => request(`/versions/${versionId}`),
  compare: (articleId, oldVerId, newVerId) => {
    return request(`/articles/${articleId}/versions/compare?old_version_id=${oldVerId}&new_version_id=${newVerId}`)
  },
  rollback: (articleId, versionId) => request(`/articles/${articleId}/versions/${versionId}/rollback`, { method: 'POST' }),
  delete: (versionId) => request(`/versions/${versionId}`, { method: 'DELETE' }),
}
