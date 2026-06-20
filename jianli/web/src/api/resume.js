import request from '@/utils/request'

export function getResumeList(userIdentity = 'default_user') {
  return request({
    url: '/resumes',
    method: 'get',
    params: { user_identity: userIdentity }
  })
}

export function getResumeDetail(id) {
  return request({
    url: `/resumes/${id}`,
    method: 'get'
  })
}

export function createResume(data) {
  return request({
    url: '/resumes',
    method: 'post',
    data
  })
}

export function updateResume(id, data) {
  return request({
    url: `/resumes/${id}`,
    method: 'put',
    data
  })
}

export function deleteResume(id) {
  return request({
    url: `/resume/${id}`,
    method: 'delete'
  })
}

export function getTemplates() {
  return request({
    url: '/templates',
    method: 'get'
  })
}

export function getTemplate(id) {
  return request({
    url: `/templates/${id}`,
    method: 'get'
  })
}

export function getVersions(resumeId) {
  return request({
    url: `/resumes/${resumeId}/versions`,
    method: 'get'
  })
}

export function getVersion(versionId) {
  return request({
    url: `/versions/${versionId}`,
    method: 'get'
  })
}

export function restoreVersion(versionId) {
  return request({
    url: `/versions/${versionId}/restore`,
    method: 'post'
  })
}

export function exportPDF(content, styleConfig, filename) {
  return request({
    url: '/pdf/export',
    method: 'post',
    data: { content, style_config: styleConfig, filename },
    responseType: 'blob'
  })
}
