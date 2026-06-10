
import axios from 'axios'

const request = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

export const getPrizes = () =&gt; request.get('/prizes')
export const doLottery = (data) =&gt; request.post('/lottery/draw', data)
export const getMyRecords = (phone) =&gt; request.get('/lottery/records', { params: { phone } })
export const claimRecord = (id) =&gt; request.put(`/lottery/records/${id}/claim`)

export default request
