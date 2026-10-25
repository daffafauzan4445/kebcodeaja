import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/auth-context'
import ProtectedRoute from './routes/protected-route'
import DashboardLayout from './components/layout/dashboard-layout'
import ComingSoon from './components/common/coming-soon'
import LoginPage from './pages/auth/login-page'
import RegisterPage from './pages/auth/register-page'
import DashboardSiswa from './pages/student/dashboard-siswa'
import DashboardGuru from './pages/teacher/dashboard-guru'

function RootRedirect() {
  const { user } = useAuth()
  if (!user) return <Navigate to="/login" replace />
  return <Navigate to={user.role === 'teacher' ? '/guru' : '/siswa'} replace />
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<RootRedirect />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route element={<ProtectedRoute allowedRoles={['student']} />}>
            <Route path="/siswa" element={<DashboardLayout />}>
              <Route index element={<DashboardSiswa />} />
              <Route path="kelas" element={<ComingSoon title="Kelas Saya" />} />
              <Route path="materi" element={<ComingSoon title="Materi" />} />
              <Route path="tantangan" element={<ComingSoon title="Tantangan" />} />
              <Route path="progres" element={<ComingSoon title="Progres Belajar" />} />
            </Route>
          </Route>

          <Route element={<ProtectedRoute allowedRoles={['teacher']} />}>
            <Route path="/guru" element={<DashboardLayout />}>
              <Route index element={<DashboardGuru />} />
              <Route path="kelas" element={<ComingSoon title="Kelola Kelas" />} />
              <Route path="materi" element={<ComingSoon title="Kelola Materi" />} />
              <Route path="tantangan" element={<ComingSoon title="Kelola Tantangan" />} />
              <Route path="rekap" element={<ComingSoon title="Rekap Nilai" />} />
            </Route>
          </Route>

          <Route
            path="*"
            element={<div className="p-8 text-slate-700">404 - Halaman tidak ditemukan</div>}
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}
