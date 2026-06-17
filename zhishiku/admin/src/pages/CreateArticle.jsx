import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { articleApi, categoryApi } from '../api'
import { useAuth } from '../context/AuthContext'
import MarkdownEditor from '../components/MarkdownEditor'
import './ArticleForm.css'

export default function CreateArticle() {
  const navigate = useNavigate()
  const { isLoggedIn } = useAuth()
  const [categories, setCategories] = useState([])
  const [form, setForm] = useState({
    title: '',
    content: '',
    category_id: '',
  })
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (!isLoggedIn) {
      navigate('/login')
      return
    }
    loadCategories()
  }, [isLoggedIn])

  const loadCategories = async () => {
    try {
      const data = await categoryApi.getAll()
      setCategories(data)
      if (data.length > 0) {
        setForm((prev) => ({ ...prev, category_id: String(data[0].id) }))
      }
    } catch (err) {
      console.error('Failed to load categories:', err)
    }
  }

  const handleSubmit = async (status) => {
    if (!form.title.trim() || !form.content.trim() || !form.category_id) {
      alert('请填写所有必填字段')
      return
    }
    setSubmitting(true)
    try {
      await articleApi.create({
        title: form.title,
        content: form.content,
        category_id: Number(form.category_id),
        status,
      })
      navigate('/')
    } catch (err) {
      alert(err.message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="article-form-page">
      <h1 className="page-title">新建文章</h1>
      <form className="article-form" onSubmit={(e) => e.preventDefault()}>
        <div className="form-group">
          <label>标题 *</label>
          <input
            type="text"
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            placeholder="请输入文章标题"
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
          <label>正文 * <span className="form-hint">（支持 Markdown 语法）</span></label>
          <MarkdownEditor
            value={form.content}
            onChange={(val) => setForm({ ...form, content: val })}
            placeholder="请输入 Markdown 内容..."
          />
        </div>

        <div className="form-actions">
          <button
            type="button"
            className="btn-primary"
            disabled={submitting}
            onClick={() => handleSubmit('published')}
          >
            {submitting ? '提交中...' : '发布文章'}
          </button>
          <button
            type="button"
            className="btn-secondary"
            disabled={submitting}
            onClick={() => handleSubmit('draft')}
          >
            {submitting ? '提交中...' : '存为草稿'}
          </button>
          <button type="button" className="btn-cancel" onClick={() => navigate(-1)}>
            取消
          </button>
        </div>
      </form>
    </div>
  )
}
