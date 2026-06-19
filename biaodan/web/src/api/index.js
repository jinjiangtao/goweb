import request from '@/utils/request'

export function listTemplates(params) {
  return request.get('/api/templates', { params })
}

export function getTemplate(id) {
  return request.get(`/api/templates/${id}`)
}

export function createTemplate(data) {
  return request.post('/api/templates', data)
}

export function updateTemplate(id, data) {
  return request.put(`/api/templates/${id}`, data)
}

export function deleteTemplate(id) {
  return request.delete(`/api/templates/${id}`)
}

export function submitForm(data) {
  return request.post('/api/submissions', data)
}

export function listSubmissions(params) {
  return request.get('/api/submissions', { params })
}

export function deleteSubmission(id) {
  return request.delete(`/api/submissions/${id}`)
}

export function exportSubmissions(params) {
  return request.get('/api/submissions/export', {
    params,
    responseType: 'blob'
  })
}

export function uploadFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/api/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function getTemplateStats(templateId) {
  return request.get(`/api/stats/${templateId}`)
}
