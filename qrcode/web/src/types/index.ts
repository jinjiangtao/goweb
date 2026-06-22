export interface StyleConfig {
  size: number
  foreground: string
  background: string
  cornerStyle: 'square' | 'rounded' | 'dot'
  eccLevel: 'L' | 'M' | 'Q' | 'H'
  margin: number
  logo: string
}

export interface Record {
  id: number
  content: string
  contentType: string
  source: string
  styleConfig: StyleConfig
  filePath: string
  status: 'active' | 'invalid'
  templateId?: number | null
  batchId?: string
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface Template {
  id: number
  name: string
  styleConfig: StyleConfig
  createdAt: string
  updatedAt: string
}

export interface ParseResult {
  raw: string
  type: string
  content: string
}

export interface Stats {
  total: number
  today: number
  templates: number
  invalid: number
  recent: Record[]
}

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export function defaultStyle(): StyleConfig {
  return {
    size: 512,
    foreground: '#0E0E11',
    background: '#FFFFFF',
    cornerStyle: 'square',
    eccLevel: 'M',
    margin: 4,
    logo: '',
  }
}
