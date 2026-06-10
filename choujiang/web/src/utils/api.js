import axios from 'axios'

const request = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

export const getPrizes = () => request.get('/prizes')
export const doLottery = (data) => request.post('/lottery/draw', data)
export const getMyRecords = (phone) => request.get('/lottery/records', { params: { phone } })
export const claimRecord = (id) => request.put(`/lottery/records/${id}/claim`)
export const submitAddress = (id, data) => request.post(`/lottery/records/${id}/address`, data)

export const getProvinces = () => request.get('/address/provinces')
export const getCities = (provinceId) => request.get('/address/cities', { params: { provinceId } })
export const getDistricts = (cityId) => request.get('/address/districts', { params: { cityId } })

export default request
