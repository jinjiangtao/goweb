import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { versionApi } from '../api'
import { useAuth } from '../context/AuthContext'
import './VersionHistory.css'

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight(str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try { return hljs.highlight(str, { language: lang }).value } catch { /* ignore */ }
    }
    try { return hljs.highlightAuto(str).value } catch { /* ignore */ }
    return ''
  },
})

export default function VersionHistory({ articleId }) {
  const navigate = useNavigate()
  const { isAdmin } = useAuth()
  const [versions, setVersions] = useState([])
  const [loading, setLoading] = useState(true)
  const [collapsed, setCollapsed] = useState(true)
  const [selectedVersions, setSelectedVersions] = useState([])
  const [detailVersion, setDetailVersion] = useState(null)
  const [detailHtml, setDetailHtml] = useState('')

  useEffect(() => {
    loadVersions()
  }, [articleId])

  const loadVersions = async () => {
    try {
      const data = await versionApi.getList(articleId)
      setVersions(data.list || [])
    } catch (err) {
      console.error('Failed to load versions:', err)
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateStr) => {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('zh-CN', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit', second: '2-digit',
    })
  }

  const handleSelectVersion = (versionId) => {
    setSelectedVersions((prev) => {
      if (prev.includes(versionId)) {
        return prev.filter((id) => id !== versionId)
      }
      if (prev.length >= 2) {
        return [prev[1], versionId]
      }
      return [...prev, versionId]
    })
  }

  const handleCompare = () => {
    if (selectedVersions.length !== 2) {
      alert('请选择两个版本进行对比')
      return
    }
    navigate(`/compare/${articleId}?old=${selectedVersions[0]}&new=${selectedVersions[1]}`)
  }

  const handleViewDetail = async (version) => {
    try {
      const data = await versionApi.getDetail(version.id)
      setDetailVersion(data)
      setDetailHtml(md.render(data.content_snapshot || ''))
    } catch (err) {
      alert(err.message)
    }
  }

  const handleRollback = async (version) => {
    if (!window.confirm(`确定要回滚到 V${version.version_number} 吗？此操作将生成一个新版本。`)) return
    try {
      await versionApi.rollback(articleId, version.id)
      alert('回滚成功')
      setDetailVersion(null)
      loadVersions()
    } catch (err) {
      alert(err.message)
    }
  }

  const handleDeleteVersion = async (version) => {
    if (!window.confirm(`确定要删除 V${version.version_number} 吗？`)) return
    try {
      await versionApi.delete(version.id)
      loadVersions()
    } catch (err) {
      alert(err.message)
    }
  }

  if (loading) return null
  if (versions.length === 0) return null

  return (
    <div className="version-history">
      <div className="version-header" onClick={() => setCollapsed(!collapsed)}>
        <h3>版本历史 ({versions.length})</h3>
        <span className="version-toggle">{collapsed ? '展开' : '收起'}</span>
      </div>

      {!collapsed && (
        <div className="version-body">
          <div className="version-actions-bar">
            <button
              className="btn-compare"
              disabled={selectedVersions.length !== 2}
              onClick={handleCompare}
            >
              对比选中版本 ({selectedVersions.length}/2)
            </button>
            {selectedVersions.length === 2 && (
              <span className="compare-hint">已选 V{versions.find(v => v.id === selectedVersions[0])?.version_number} vs V{versions.find(v => v.id === selectedVersions[1])?.version_number}</span>
            )}
          </div>

          <div className="version-list">
            {versions.map((v) => (
              <div
                key={v.id}
                className={`version-item ${selectedVersions.includes(v.id) ? 'selected' : ''}`}
              >
                <div className="version-item-left">
                  <input
                    type="checkbox"
                    checked={selectedVersions.includes(v.id)}
                    onChange={() => handleSelectVersion(v.id)}
                  />
                  <span className="version-number">V{v.version_number}</span>
                  <span className={`version-status ${v.status_snapshot}`}>
                    {v.status_snapshot === 'published' ? '已发布' : '草稿'}
                  </span>
                  {v.is_rollback && (
                    <span className="version-rollback-badge">回滚</span>
                  )}
                </div>
                <div className="version-item-right">
                  <span className="version-time">{formatDate(v.created_at)}</span>
                  <button className="btn-view-ver" onClick={() => handleViewDetail(v)}>查看</button>
                  {isAdmin && (
                    <button className="btn-del-ver" onClick={() => handleDeleteVersion(v)}>删除</button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {detailVersion && (
        <div className="modal-overlay" onClick={() => setDetailVersion(null)}>
          <div className="modal-content version-detail-modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>版本 V{detailVersion.version_number} 详情</h3>
              <button className="modal-close" onClick={() => setDetailVersion(null)}>×</button>
            </div>
            <div className="modal-body">
              <div className="version-meta-info">
                <span>保存时间：{formatDate(detailVersion.created_at)}</span>
                <span className={`version-status ${detailVersion.status_snapshot}`}>
                  {detailVersion.status_snapshot === 'published' ? '已发布' : '草稿'}
                </span>
                {detailVersion.is_rollback && detailVersion.rollback_from_ver && (
                  <span className="version-rollback-badge">回滚自 V{detailVersion.rollback_from_ver}</span>
                )}
              </div>
              <h4 className="detail-title">{detailVersion.title_snapshot}</h4>
              <div
                className="detail-content markdown-body"
                dangerouslySetInnerHTML={{ __html: detailHtml }}
              />
            </div>
            <div className="modal-footer">
              <button
                className="btn-rollback"
                onClick={() => handleRollback(detailVersion)}
              >
                回滚到此版本
              </button>
              <button className="btn-close-modal" onClick={() => setDetailVersion(null)}>关闭</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
