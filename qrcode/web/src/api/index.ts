import client from './client'
import type { ApiResponse, StyleConfig, Record, Template, ParseResult, Stats } from '@/types'

export async function previewQR(content: string, style: StyleConfig): Promise<Blob> {
  const res = await client.post('/qrcode/preview', { content, styleConfig: style }, { responseType: 'blob' })
  return res.data
}

export async function generateQR(content: string, style: StyleConfig, opts?: { saveTemplate?: boolean; templateName?: string; templateId?: number }) {
  const res = await client.post<ApiResponse<{ id: number; previewUrl: string; record: Record }>>('/qrcode/generate', {
    content,
    styleConfig: style,
    ...opts,
  })
  return res.data.data
}

export async function batchGenerate(items: { content: string; styleConfig?: StyleConfig }[], style: StyleConfig) {
  const res = await client.post<ApiResponse<{ batchId: string; count: number; records: Record[] }>>('/qrcode/batch', {
    items,
    styleConfig: style,
  })
  return res.data.data
}

export async function parseText(texts: string[]) {
  const res = await client.post<ApiResponse<{ results: ParseResult[] }>>('/parse/text', { texts })
  return res.data.data.results
}

export async function parseImage(files: File[]) {
  const form = new FormData()
  files.forEach((f) => form.append('files', f))
  const res = await client.post<ApiResponse<{ results: ParseResult[] }>>('/parse/image', form)
  return res.data.data.results
}

export async function listTemplates(keyword?: string) {
  const res = await client.get<ApiResponse<Template[]>>('/templates', { params: { keyword } })
  return res.data.data
}

export async function createTemplate(name: string, styleConfig: StyleConfig) {
  const res = await client.post<ApiResponse<Template>>('/templates', { name, styleConfig })
  return res.data.data
}

export async function updateTemplate(id: number, name: string, styleConfig: StyleConfig) {
  const res = await client.put<ApiResponse<Template>>(`/templates/${id}`, { name, styleConfig })
  return res.data.data
}

export async function deleteTemplate(id: number) {
  await client.delete(`/templates/${id}`)
}

export async function listRecords(params?: { status?: string; keyword?: string; page?: number; size?: number }) {
  const res = await client.get<ApiResponse<{ total: number; list: Record[] }>>('/records', { params })
  return res.data.data
}

export async function getRecord(id: number) {
  const res = await client.get<ApiResponse<Record>>(`/records/${id}`)
  return res.data.data
}

export async function updateRecordStatus(id: number, status: string) {
  await client.patch(`/records/${id}/status`, { status })
}

export async function updateRecordContent(id: number, content: string) {
  const res = await client.put<ApiResponse<Record>>(`/records/${id}/content`, { content })
  return res.data.data
}

export async function deleteRecord(id: number) {
  await client.delete(`/records/${id}`)
}

export async function exportRecords(ids: number[]) {
  const res = await client.post('/records/export', { ids }, { responseType: 'blob' })
  return res.data as Blob
}

export async function getStats() {
  const res = await client.get<ApiResponse<Stats>>('/stats')
  return res.data.data
}

export function recordImageUrl(id: number) {
  return `/api/records/${id}/image`
}

export function templatePreviewUrl(id: number) {
  return `/api/templates/${id}/preview`
}
