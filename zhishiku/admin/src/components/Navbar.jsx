import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { articleApi } from '../api'
import './Navbar.css'

export default function Navbar({ onSearch }) {
  const { user, logout, isAdmin } = useAuth()
  const [keyword, setSearchKeyword] = useState('')
  const navigate = useNavigate()

  const handleSearch = (e) => {
    e.preventDefault()
    if (onSearch) {
      onSearch(keyword)
    } else {
      navigate(`/?keyword=${encodeURIComponent(keyword)}`)
    }
  }

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  return (
    <nav className="navbar">
      <div className="navbar-inner">
        <Link to="/" className="navbar-brand">知识库</Link>

        <form className="navbar-search" onSubmit={handleSearch}>
          <input
            type="text"
            placeholder="搜索文章..."
            value={keyword}
            onChange={(e) => setSearchKeyword(e.target.value)}
            className="search-input"
          />
          <button type="submit" className="search-btn">搜索</button>
        </form>

        <div className="navbar-right">
          {user ? (
            <>
              <span className="user-nickname">{user.nickname}</span>
              <Link to="/my-articles" className="nav-link">我的文章</Link>
              {isAdmin && <Link to="/categories" className="nav-link">分类管理</Link>}
              <Link to="/article/create" className="nav-link btn-create">新建文章</Link>
              <button onClick={handleLogout} className="btn-logout">退出</button>
            </>
          ) : (
            <>
              <Link to="/login" className="nav-link">登录</Link>
              <Link to="/register" className="nav-link">注册</Link>
            </>
          )}
        </div>
      </div>
    </nav>
  )
}
