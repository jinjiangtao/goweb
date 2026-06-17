import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { articleApi } from '../api'
import { useAuth } from '../context/AuthContext'
import './MyArticles.css'

export default function MyArticles() {
  const { user, isAdmin } = useAuth()
  const navigate = useNavigate()
  const [articles, setArticles] = useState([])
  const [totalPages, setTotalPages] = useState(1)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      navigate('/login')
      return
    }
    loadArticles()
  }, [page, user])

  const loadArticles = async () => {
    setLoading(true)
    try {
      const data = await articleApi.getMine({ page, page_size: 10 })
      setArticles(data.list || [])
      setTotalPages(data.total_pages || 1)
    } catch (err) {
      console.error('Failed to load articles:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('确定要删除这篇文章吗？')) return
    try {
      await articleApi.delete(id)
      loadArticles()
    } catch (err) {
      alert(err.message)
    }
  }

  const formatDate = (dateStr) => {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('zh-CN', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit',
    })
  }

  return (
    <div className="my-articles-page">
      <div className="page-header">
        <h1 className="page-title">我的文章</h1>
        <Link to="/article/create" className="btn-primary">新建文章</Link>
      </div>

      {loading ? (
        <div className="loading">加载中...</div>
      ) : articles.length === 0 ? (
        <div className="empty">暂无文章，快去创建一篇吧！</div>
      ) : (
        <>
          <div className="my-article-list">
            {articles.map((article) => (
              <div key={article.id} className="my-article-card">
                <div className="my-article-info">
                  <Link to={`/article/${article.id}`} className="my-article-title">
                    {article.title}
                  </Link>
                  <div className="my-article-meta">
                    <span className="article-category">{article.category_name}</span>
                    <span className={`status-badge ${article.status}`}>
                      {article.status === 'draft' ? '草稿' : '已发布'}
                    </span>
                    <span className="article-date">{formatDate(article.created_at)}</span>
                  </div>
                </div>
                <div className="my-article-actions">
                  <Link to={`/article/edit/${article.id}`} className="btn-sm btn-edit">编辑</Link>
                  <button onClick={() => handleDelete(article.id)} className="btn-sm btn-del">删除</button>
                </div>
              </div>
            ))}
          </div>

          {totalPages > 1 && (
            <div className="pagination">
              <button disabled={page <= 1} onClick={() => setPage(page - 1)} className="page-btn">
                上一页
              </button>
              <span className="page-info">{page} / {totalPages}</span>
              <button disabled={page >= totalPages} onClick={() => setPage(page + 1)} className="page-btn">
                下一页
              </button>
            </div>
          )}
        </>
      )}
    </div>
  )
}
