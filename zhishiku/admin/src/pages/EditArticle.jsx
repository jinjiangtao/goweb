import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { articleApi, categoryApi } from '../api'
import { useAuth } from '../context/AuthContext'
import './ArticleForm.css'

export default function EditArticle() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user, isAdmin } = useAuth()
  const [categories, setCategories] = useState([])
  const [form, setForm] = useState({
    title: '',
    content: '',
    category_id: '',
    status: 'draft',
  })
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    loadData()
  }, [id])

  const loadData = async () => {
    try {
      const [articleData, categoriesData] = await Promise.all([
        articleApi.getDetail(id),
        categoryApi.getAll(),
      ])

      if (articleData.author_id !== user?.id && !isAdmin) {
        alert('无权编辑此文章')
        navigate('/')
        return
      }

      setForm({
        title: articleData.title,
        content: articleData.content,
        category_id: String(articleData.category_id),
        status: articleData.status,
      })
      setCategories(categoriesData)
    } catch (err) {
      alert(err.message)
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!form.title.trim() || !form.content.trim() || !form.category_id) {
      alert('请填写所有必填字段')
      return
    }
    setSubmitting(true)
    try {
      await articleApi.update(id, {
        title: form.title,
        content: form.content,
        category_id: Number(form.category_id),
        status: form.status,
      })
      navigate('/')
    } catch (err) {
      alert(err.message)
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <div className="loading">加载中...</div>

  return (
    <div className="article-form-page">
      <h1 className="page-title">编辑文章</h1>
      <form onSubmit={handleSubmit} className="article-form">
        <div className="form-group">
          <label>标题 *</label>
          <input
            type="text"
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            required
          />
        </div>

        <div className="form-group">
          <label>分类 *</label>
          <select
            value={form.category_id}
            onChange={(e) => setForm({ ...form, category_id: e.target.value })}
            required
          >
            <option value="">请选择分类</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>{cat.name}</option>
            ))}
          </select>
        </div>

        <div className="form-group">
          <label>正文 *</label>
          <textarea
            value={form.content}
            onChange={(e) => setForm({ ...form, content: e.target.value })}
            rows={15}
            required
          />
        </div>

        <div className="form-group">
          <label>状态</label>
          <select
            value={form.status}
            onChange={(e) => setForm({ ...form, status: e.target.value })}
          >
            <option value="draft">草稿</option>
            <option value="published">已发布</option>
          </select>
        </div>

        <div className="form-actions">
          <button type="submit" className="btn-primary" disabled={submitting}>
            {submitting ? '提交中...' : '保存修改'}
          </button>
          <button type="button" className="btn-cancel" onClick={() => navigate(-1)}>
            取消
          </button>
        </div>
      </form>
    </div>
  )
}
