import axios from 'axios'

const request = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

export const getPrizes = () => request.get('/prizes')
export const doLottery = (data) => request.post('/lottery/draw', data)
export const getMyRecords = (phone) => request.get('/lottery/records', { params: { phone } })
export const claimRecord = (id) => request.put(`/lottery/records/${id}/claim`)

export default request
