import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'

vi.mock('axios')

const mockAxios = axios as vi.Mocked<typeof axios>

describe('API Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockAxios.create.mockReturnValue({
      post: vi.fn(),
      get: vi.fn()
    } as any)
  })

  it('login should return token', async () => {
    const apiModule = await import('@/api')
    const mockPost = vi.fn().mockResolvedValue({ data: { token: 'test-token' } })
    mockAxios.create.mockReturnValue({ post: mockPost } as any)
    
    const result = await apiModule.login({ username: 'test', password: 'test' })
    expect(result.token).toBe('test-token')
  })

  it('register should return user', async () => {
    const apiModule = await import('@/api')
    const mockPost = vi.fn().mockResolvedValue({ data: { id: '1', username: 'test', nickname: 'Test', avatar: '' } })
    mockAxios.create.mockReturnValue({ post: mockPost } as any)
    
    const result = await apiModule.register({ username: 'test', password: 'test', nickname: 'Test' })
    expect(result.username).toBe('test')
  })

  it('getFriends should return friends list', async () => {
    const apiModule = await import('@/api')
    const mockGet = vi.fn().mockResolvedValue({ data: [{ id: '2', username: 'friend', nickname: 'Friend', avatar: '', online: true }] })
    mockAxios.create.mockReturnValue({ get: mockGet } as any)
    
    const result = await apiModule.getFriends('1')
    expect(result.length).toBe(1)
    expect(result[0].nickname).toBe('Friend')
  })
})