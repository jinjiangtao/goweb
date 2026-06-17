import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { articleApi } from '../api'
import { useAuth } from '../context/AuthContext'
import './ArticleDetail.css'

export default function ArticleDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user, isAdmin } = useAuth()
  const [article, setArticle] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadArticle()
  }, [id])

  const loadArticle = async () => {
    try {
      const data = await articleApi.getDetail(id)
      setArticle(data)
    } catch (err) {
      console.error('Failed to load article:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!window.confirm('确定要删除这篇文章吗？')) return
    try {
      await articleApi.delete(id)
      navigate('/')
    } catch (err) {
      alert(err.message)
    }
  }

  const canEdit = user && (article?.author_id === user.id || isAdmin)

  const formatDate = (dateStr) => {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('zh-CN', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit',
    })
  }

  if (loading) return <div className="loading">加载中...</div>
  if (!article) return <div className="loading">文章不存在</div>

  return (
    <div className="article-detail-page">
      <article className="article-detail">
        <header className="article-header">
          <h1 className="article-title">{article.title}</h1>
          <div className="article-meta">
            <span className="meta-item category-badge">{article.category_name}</span>
            <span className="meta-item">{article.author_name}</span>
            <span className="meta-item">{formatDate(article.published_at || article.created_at)}</span>
          </div>
        </header>

        <div className="article-content">
          {article.content.split('\n').map((paragraph, idx) => (
            <p key={idx}>{paragraph}</p>
          ))}
        </div>

        {canEdit && (
          <div className="article-actions">
            <Link to={`/article/edit/${article.id}`} className="btn-edit">编辑</Link>
            <button onClick={handleDelete} className="btn-delete">删除</button>
          </div>
        )}
      </article>
    </div>
  )
}
