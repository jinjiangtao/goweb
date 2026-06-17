import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import Navbar from './components/Navbar'
import ArticleList from './pages/ArticleList'
import ArticleDetail from './pages/ArticleDetail'
import CreateArticle from './pages/CreateArticle'
import EditArticle from './pages/EditArticle'
import MyArticles from './pages/MyArticles'
import Categories from './pages/Categories'
import Login from './pages/Login'
import Register from './pages/Register'
import VersionCompare from './pages/VersionCompare'

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Navbar />
        <Routes>
          <Route path="/" element={<ArticleList />} />
          <Route path="/article/:id" element={<ArticleDetail />} />
          <Route path="/article/create" element={<CreateArticle />} />
          <Route path="/article/edit/:id" element={<EditArticle />} />
          <Route path="/my-articles" element={<MyArticles />} />
          <Route path="/categories" element={<Categories />} />
          <Route path="/compare/:id" element={<VersionCompare />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
