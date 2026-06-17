import { useState, useEffect } from 'react'
import { categoryApi } from '../api'
import { useAuth } from '../context/AuthContext'
import './Categories.css'

export default function Categories() {
  const { isAdmin } = useAuth()
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(true)
  const [editingId, setEditingId] = useState(null)
  const [form, setForm] = useState({ name: '', sort_order: 0 })
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    loadCategories()
  }, [])

  const loadCategories = async () => {
    try {
      const data = await categoryApi.getAll()
      setCategories(data)
    } catch (err) {
      console.error('Failed to load categories:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (e) => {
    e.preventDefault()
    if (!form.name.trim()) return
    try {
      await categoryApi.create(form)
      setForm({ name: '', sort_order: 0 })
      setShowForm(false)
      loadCategories()
    } catch (err) {
      alert(err.message)
    }
  }

  const handleUpdate = async (e) => {
    e.preventDefault()
    if (!form.name.trim()) return
    try {
      await categoryApi.update(editingId, form)
      setEditingId(null)
      setForm({ name: '', sort_order: 0 })
      loadCategories()
    } catch (err) {
      alert(err.message)
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('确定要删除此分类吗？')) return
    try {
      await categoryApi.delete(id)
      loadCategories()
    } catch (err) {
      alert(err.message)
    }
  }

  const startEdit = (cat) => {
    setEditingId(cat.id)
    setForm({ name: cat.name, sort_order: cat.sort_order })
    setShowForm(false)
  }

  const cancelEdit = () => {
    setEditingId(null)
    setForm({ name: '', sort_order: 0 })
  }

  if (!isAdmin) {
    return (
      <div className="categories-page">
        <div className="no-permission">需要管理员权限才能访问此页面</div>
      </div>
    )
  }

  return (
    <div className="categories-page">
      <div className="page-header">
        <h1 className="page-title">分类管理</h1>
        <button
          className="btn-primary"
          onClick={() => { setShowForm(!showForm); setEditingId(null); setForm({ name: '', sort_order: 0 }) }}
        >
          {showForm ? '取消' : '新增分类'}
        </button>
      </div>

      {showForm && (
        <form className="category-form" onSubmit={handleCreate}>
          <div className="form-group">
            <label>分类名称 *</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
            />
          </div>
          <div className="form-group">
            <label>排序号</label>
            <input
              type="number"
              value={form.sort_order}
              onChange={(e) => setForm({ ...form, sort_order: Number(e.target.value) })}
            />
          </div>
          <button type="submit" className="btn-primary">创建</button>
        </form>
      )}

      {loading ? (
        <div className="loading">加载中...</div>
      ) : (
        <div className="category-list">
          {categories.map((cat) => (
            <div key={cat.id} className="category-card">
              {editingId === cat.id ? (
                <form className="inline-form" onSubmit={handleUpdate}>
                  <input
                    type="text"
                    value={form.name}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                    required
                  />
                  <input
                    type="number"
                    value={form.sort_order}
                    onChange={(e) => setForm({ ...form, sort_order: Number(e.target.value) })}
                    style={{ width: '80px' }}
                  />
                  <button type="submit" className="btn-sm btn-save">保存</button>
                  <button type="button" className="btn-sm btn-cancel-sm" onClick={cancelEdit}>取消</button>
                </form>
              ) : (
                <>
                  <div className="category-info">
                    <span className="category-name">{cat.name}</span>
                    <span className="category-sort">排序: {cat.sort_order}</span>
                  </div>
                  <div className="category-actions">
                    <button className="btn-sm btn-edit" onClick={() => startEdit(cat)}>编辑</button>
                    <button className="btn-sm btn-del" onClick={() => handleDelete(cat.id)}>删除</button>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
