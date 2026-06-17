import { useState, useEffect, useMemo } from 'react'
import { useParams, useSearchParams, useNavigate } from 'react-router-dom'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { versionApi } from '../api'
import './VersionCompare.css'

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

export default function VersionCompare() {
  const { id } = useParams()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const oldVerId = searchParams.get('old')
  const newVerId = searchParams.get('new')

  const [diffData, setDiffData] = useState(null)
  const [oldVersion, setOldVersion] = useState(null)
  const [newVersion, setNewVersion] = useState(null)
  const [loading, setLoading] = useState(true)
  const [viewMode, setViewMode] = useState('inline')
  const [renderMode, setRenderMode] = useState('text')

  useEffect(() => {
    if (oldVerId && newVerId) {
      loadData()
    }
  }, [id, oldVerId, newVerId])

  const loadData = async () => {
    setLoading(true)
    try {
      const [diffResult, oldDetail, newDetail] = await Promise.all([
        versionApi.compare(id, oldVerId, newVerId),
        versionApi.getDetail(oldVerId),
        versionApi.getDetail(newVerId),
      ])
      setDiffData(diffResult)
      setOldVersion(oldDetail)
      setNewVersion(newDetail)
    } catch (err) {
      alert(err.message)
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

  const getLineClass = (type) => {
    switch (type) {
      case 'added': return 'diff-added'
      case 'removed': return 'diff-removed'
      case 'unchanged': return 'diff-unchanged'
      default: return ''
    }
  }

  const getLinePrefix = (type) => {
    switch (type) {
      case 'added': return '+'
      case 'removed': return '-'
      default: return ' '
    }
  }

  const renderMarkdownLine = (content) => {
    const html = md.render(content)
    return <span dangerouslySetInnerHTML={{ __html: html }} />
  }

  const oldHtml = useMemo(() => oldVersion ? md.render(oldVersion.content_snapshot || '') : '', [oldVersion])
  const newHtml = useMemo(() => newVersion ? md.render(newVersion.content_snapshot || '') : '', [newVersion])

  if (loading) return <div className="loading">加载中...</div>

  if (!diffData) {
    return (
      <div className="compare-page">
        <div className="compare-empty">请从文章详情页选择两个版本进行对比</div>
      </div>
    )
  }

  return (
    <div className="compare-page">
      <div className="compare-header">
        <h1>版本对比</h1>
        <button className="btn-back" onClick={() => navigate(`/article/${id}`)}>返回文章</button>
      </div>

      <div className="compare-controls">
        <div className="control-group">
          <label>展示模式：</label>
          <select value={viewMode} onChange={(e) => setViewMode(e.target.value)}>
            <option value="inline">逐行对比</option>
            <option value="side">左右对比</option>
          </select>
        </div>
        <div className="control-group">
          <label>渲染方式：</label>
          <select value={renderMode} onChange={(e) => setRenderMode(e.target.value)}>
            <option value="text">纯文本</option>
            <option value="markdown">Markdown 渲染</option>
          </select>
        </div>
      </div>

      {(diffData.title_diff?.changed || diffData.category_diff?.changed) && (
        <div className="compare-summary">
          {diffData.title_diff?.changed && (
            <div className="summary-item">
              <span className="summary-label">标题：</span>
              <span className="summary-old">{diffData.title_diff.old_value}</span>
              <span className="summary-arrow">→</span>
              <span className="summary-new">{diffData.title_diff.new_value}</span>
            </div>
          )}
          {diffData.category_diff?.changed && (
            <div className="summary-item">
              <span className="summary-label">分类：</span>
              <span className="summary-old">{diffData.category_diff.old_name}</span>
              <span className="summary-arrow">→</span>
              <span className="summary-new">{diffData.category_diff.new_name}</span>
            </div>
          )}
        </div>
      )}

      {!diffData.has_change && (
        <div className="compare-no-diff">两个版本内容完全一致</div>
      )}

      {diffData.has_change && viewMode === 'inline' && renderMode === 'text' && (
        <div className="compare-inline">
          <div className="diff-legend">
            <span className="legend-added">+ 新增</span>
            <span className="legend-removed">- 删除</span>
            <span className="legend-unchanged">&nbsp; 未变</span>
          </div>
          <div className="diff-table">
            {diffData.lines.map((line, idx) => (
              <div key={idx} className={`diff-line ${getLineClass(line.type)}`}>
                <span className="diff-prefix">{getLinePrefix(line.type)}</span>
                <span className="diff-line-num">
                  {line.old_line || ''}
                </span>
                <span className="diff-line-num">
                  {line.new_line || ''}
                </span>
                <span className="diff-content">{line.content}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {diffData.has_change && viewMode === 'inline' && renderMode === 'markdown' && (
        <div className="compare-inline">
          <div className="diff-legend">
            <span className="legend-added">+ 新增</span>
            <span className="legend-removed">- 删除</span>
            <span className="legend-unchanged">&nbsp; 未变</span>
          </div>
          <div className="diff-table">
            {diffData.lines.map((line, idx) => (
              <div key={idx} className={`diff-line ${getLineClass(line.type)}`}>
                <span className="diff-prefix">{getLinePrefix(line.type)}</span>
                <span className="diff-line-num">
                  {line.old_line || ''}
                </span>
                <span className="diff-line-num">
                  {line.new_line || ''}
                </span>
                <span className="diff-content markdown-body">
                  {renderMarkdownLine(line.content, line.type)}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {diffData.has_change && viewMode === 'side' && renderMode === 'text' && (
        <div className="compare-side">
          <div className="side-pane">
            <div className="side-header">
              <span>V{oldVersion?.version_number}</span>
              <span className="side-time">{formatDate(oldVersion?.created_at)}</span>
            </div>
            <div className="side-content">
              {diffData.lines.filter(l => l.type === 'removed' || l.type === 'unchanged').map((line, idx) => (
                <div key={idx} className={`diff-line ${getLineClass(line.type)}`}>
                  <span className="diff-prefix">{getLinePrefix(line.type)}</span>
                  <span className="diff-content">{line.content}</span>
                </div>
              ))}
            </div>
          </div>
          <div className="side-pane">
            <div className="side-header">
              <span>V{newVersion?.version_number}</span>
              <span className="side-time">{formatDate(newVersion?.created_at)}</span>
            </div>
            <div className="side-content">
              {diffData.lines.filter(l => l.type === 'added' || l.type === 'unchanged').map((line, idx) => (
                <div key={idx} className={`diff-line ${getLineClass(line.type)}`}>
                  <span className="diff-prefix">{getLinePrefix(line.type)}</span>
                  <span className="diff-content">{line.content}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {diffData.has_change && viewMode === 'side' && renderMode === 'markdown' && (
        <div className="compare-side">
          <div className="side-pane">
            <div className="side-header">
              <span>V{oldVersion?.version_number}</span>
              <span className="side-time">{formatDate(oldVersion?.created_at)}</span>
            </div>
            <div
              className="side-content markdown-body"
              dangerouslySetInnerHTML={{ __html: oldHtml }}
            />
          </div>
          <div className="side-pane">
            <div className="side-header">
              <span>V{newVersion?.version_number}</span>
              <span className="side-time">{formatDate(newVersion?.created_at)}</span>
            </div>
            <div
              className="side-content markdown-body"
              dangerouslySetInnerHTML={{ __html: newHtml }}
            />
          </div>
        </div>
      )}
    </div>
  )
}
