import request from './request'

export const getDashboard = () => request.get('/dashboard')

export const getQuestions = (params) => request.get('/questions', { params })

export const getAllQuestions = (params) => request.get('/questions/all', { params })

export const getQuestionCategories = () => request.get('/questions/categories')

export const getQuestion = (id) => request.get(`/questions/${id}`)

export const createQuestion = (data) => request.post('/questions', data)

export const batchCreateQuestions = (data) => request.post('/questions/batch', data)

export const updateQuestion = (id, data) => request.put(`/questions/${id}`, data)

export const deleteQuestion = (id) => request.delete(`/questions/${id}`)

export const batchDeleteQuestions = (ids) => request.post('/questions/delete-batch', { ids })

export const generateExam = (data) => request.post('/exam/generate', data)

export const submitExam = (data) => request.post('/exam/submit', data)

export const getExamRecords = (params) => request.get('/exam/records', { params })

export const getExamRecord = (id) => request.get(`/exam/records/${id}`)

export const getWrongQuestions = (params) => request.get('/wrong', { params })

export const getWrongPractice = (params) => request.get('/wrong/practice', { params })

export const deleteWrongQuestion = (id, userId) => request.delete(`/wrong/${id}`, { params: { user_id: userId } })

export const clearWrongQuestions = (userId) => request.delete('/wrong', { params: { user_id: userId } })

export const removeWrongAfterCorrect = (data) => request.post('/wrong/remove-correct', data)

export const getStatistics = (params) => request.get('/stats', { params })
