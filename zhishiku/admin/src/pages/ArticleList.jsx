import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { articleApi } from '../api'
import { useAuth } from '../context/AuthContext'
import './ArticleList.css'

export default function ArticleList() {
  const [articles, setArticles] = useState([])
  const [totalPages, setTotalPages] = useState(1)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const { isLoggedIn } = useAuth()
  const [searchParams] = useSearchParams()
  const keyword = searchParams.get('keyword') || ''

  useEffect(() => {
    loadArticles()
  }, [page, keyword])

  const loadArticles = async () => {
    setLoading(true)
    try {
      const params = { page, page_size: 10 }
      if (keyword) params.keyword = keyword
      const data = await articleApi.getList(params)
      setArticles(data.list || [])
      setTotalPages(data.total_pages || 1)
    } catch (err) {
      console.error('Failed to load articles:', err)
    } finally {
      setLoading(false)
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
    <div className="article-list-page">
      <div className="page-header">
        <h1 className="page-title">{keyword ? `搜索: "${keyword}"` : '知识库文章'}</h1>
        {isLoggedIn && (
          <Link to="/article/create" className="btn-primary">新建文章</Link>
        )}
      </div>

      {loading ? (
        <div className="loading">加载中...</div>
      ) : articles.length === 0 ? (
        <div className="empty">暂无文章</div>
      ) : (
        <>
          <div className="article-list">
            {articles.map((article) => (
              <Link to={`/article/${article.id}`} key={article.id} className="article-card">
                <h2 className="article-card-title">{article.title}</h2>
                <div className="article-card-meta">
                  <span className="article-category">{article.category_name}</span>
                  <span className="article-author">{article.author_name}</span>
                  <span className="article-date">{formatDate(article.published_at || article.created_at)}</span>
                </div>
              </Link>
            ))}
          </div>

          {totalPages > 1 && (
            <div className="pagination">
              <button
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
                className="page-btn"
              >
                上一页
              </button>
              <span className="page-info">{page} / {totalPages}</span>
              <button
                disabled={page >= totalPages}
                onClick={() => setPage(page + 1)}
                className="page-btn"
              >
                下一页
              </button>
            </div>
          )}
        </>
      )}
    </div>
  )
}
