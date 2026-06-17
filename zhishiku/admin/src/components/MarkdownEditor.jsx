import { useRef, useCallback, useEffect, useMemo } from 'react'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import './MarkdownEditor.css'

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight(str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch {
        // ignore
      }
    }
    try {
      return hljs.highlightAuto(str).value
    } catch {
      // ignore
    }
    return ''
  },
})

export default function MarkdownEditor({ value, onChange, placeholder }) {
  const editRef = useRef(null)
  const previewRef = useRef(null)
  const syncingRef = useRef(false)

  const html = useMemo(() => md.render(value || ''), [value])

  const handleEditScroll = useCallback(() => {
    if (syncingRef.current) return
    syncingRef.current = true
    const edit = editRef.current
    const preview = previewRef.current
    if (edit && preview) {
      const ratio = edit.scrollTop / (edit.scrollHeight - edit.clientHeight || 1)
      preview.scrollTop = ratio * (preview.scrollHeight - preview.clientHeight || 1)
    }
    requestAnimationFrame(() => { syncingRef.current = false })
  }, [])

  const handlePreviewScroll = useCallback(() => {
    if (syncingRef.current) return
    syncingRef.current = true
    const edit = editRef.current
    const preview = previewRef.current
    if (edit && preview) {
      const ratio = preview.scrollTop / (preview.scrollHeight - preview.clientHeight || 1)
      edit.scrollTop = ratio * (edit.scrollHeight - edit.clientHeight || 1)
    }
    requestAnimationFrame(() => { syncingRef.current = false })
  }, [])

  useEffect(() => {
    const edit = editRef.current
    const preview = previewRef.current
    edit?.addEventListener('scroll', handleEditScroll)
    preview?.addEventListener('scroll', handlePreviewScroll)
    return () => {
      edit?.removeEventListener('scroll', handleEditScroll)
      preview?.removeEventListener('scroll', handlePreviewScroll)
    }
  }, [handleEditScroll, handlePreviewScroll])

  const handleKeyDown = (e) => {
    if (e.key === 'Tab') {
      e.preventDefault()
      const textarea = e.target
      const start = textarea.selectionStart
      const end = textarea.selectionEnd
      const newValue = value.substring(0, start) + '  ' + value.substring(end)
      onChange(newValue)
      setTimeout(() => {
        textarea.selectionStart = textarea.selectionEnd = start + 2
      }, 0)
    }
  }

  return (
    <div className="md-editor">
      <div className="md-editor-pane">
        <div className="md-pane-header">编辑</div>
        <textarea
          ref={editRef}
          className="md-editor-input"
          value={value || ''}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder || '请输入 Markdown 内容...'}
        />
      </div>
      <div className="md-editor-divider" />
      <div className="md-editor-pane">
        <div className="md-pane-header">预览</div>
        <div
          ref={previewRef}
          className="md-editor-preview markdown-body"
          dangerouslySetInnerHTML={{ __html: html }}
        />
      </div>
    </div>
  )
}
