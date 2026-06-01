import { describe, it, expect, vi, beforeEach } from 'vitest'
import axios from 'axios'
import * as api from '@/api'

vi.mock('axios')

describe('API Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('login should return token', async () => {
    const mockResponse = { data: { token: 'test-token' } }
    vi.mocked(axios.create).mockReturnValue({
      post: vi.fn().mockResolvedValue(mockResponse)
    } as any)

    const result = await api.login({ username: 'test', password: 'test' })
    expect(result.token).toBe('test-token')
  })

  it('register should return user', async () => {
    const mockUser = { id: '1', username: 'test', nickname: 'Test', avatar: '' }
    vi.mocked(axios.create).mockReturnValue({
      post: vi.fn().mockResolvedValue({ data: mockUser })
    } as any)

    const result = await api.register({ username: 'test', password: 'test', nickname: 'Test' })
    expect(result.username).toBe('test')
  })

  it('getFriends should return friends list', async () => {
    const mockFriends = [{ id: '2', username: 'friend', nickname: 'Friend', avatar: '', online: true }]
    vi.mocked(axios.create).mockReturnValue({
      get: vi.fn().mockResolvedValue({ data: mockFriends })
    } as any)

    const result = await api.getFriends('1')
    expect(result.length).toBe(1)
    expect(result[0].nickname).toBe('Friend')
  })
})