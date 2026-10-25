import { useState } from 'react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import {
  LayoutDashboard, BookOpen, Trophy, BarChart3, Users, FileText,
  ClipboardList, LogOut, Menu, X, Code2,
} from 'lucide-react'
import { useAuth } from '../../context/auth-context'

const menus = {
  student: [
    { to: '/siswa', label: 'Dashboard', icon: LayoutDashboard, end: true },
    { to: '/siswa/kelas', label: 'Kelas Saya', icon: Users },
    { to: '/siswa/materi', label: 'Materi', icon: BookOpen },
    { to: '/siswa/tantangan', label: 'Tantangan', icon: Trophy },
    { to: '/siswa/progres', label: 'Progres Belajar', icon: BarChart3 },
  ],
  teacher: [
    { to: '/guru', label: 'Dashboard', icon: LayoutDashboard, end: true },
    { to: '/guru/kelas', label: 'Kelola Kelas', icon: Users },
    { to: '/guru/materi', label: 'Kelola Materi', icon: FileText },
    { to: '/guru/tantangan', label: 'Kelola Tantangan', icon: ClipboardList },
    { to: '/guru/rekap', label: 'Rekap Nilai', icon: BarChart3 },
  ],
}

export default function DashboardLayout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const [open, setOpen] = useState(false)
  const items = menus[user.role]

  const handleLogout = () => {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <div className="min-h-screen bg-slate-50 md:flex">
      {open && (
        <div className="fixed inset-0 bg-black/40 z-30 md:hidden" onClick={() => setOpen(false)} />
      )}

      <aside
        className={`fixed md:sticky top-0 z-40 h-screen w-64 bg-white border-r border-slate-200 flex flex-col transition-transform ${
          open ? 'translate-x-0' : '-translate-x-full'
        } md:translate-x-0`}
      >
        <div className="flex items-center justify-between px-5 h-16 border-b border-slate-200">
          <div className="flex items-center gap-2 text-indigo-600 font-bold text-lg">
            <Code2 className="w-6 h-6" /> keBCodeAja
          </div>
          <button className="md:hidden" onClick={() => setOpen(false)}>
            <X className="w-5 h-5" />
          </button>
        </div>

        <nav className="flex-1 p-3 space-y-1 overflow-y-auto">
          {items.map(({ to, label, icon: Icon, end }) => (
            <NavLink
              key={to} to={to} end={end} onClick={() => setOpen(false)}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition ${
                  isActive ? 'bg-indigo-50 text-indigo-700' : 'text-slate-600 hover:bg-slate-100'
                }`
              }
            >
              <Icon className="w-5 h-5" /> {label}
            </NavLink>
          ))}
        </nav>

        <div className="p-3 border-t border-slate-200">
          <div className="px-3 py-2 mb-1">
            <p className="text-sm font-semibold text-slate-800 truncate">{user.name}</p>
            <p className="text-xs text-slate-500">{user.role === 'teacher' ? 'Guru' : 'Siswa'}</p>
          </div>
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-red-600 hover:bg-red-50"
          >
            <LogOut className="w-5 h-5" /> Keluar
          </button>
        </div>
      </aside>

      <div className="flex-1 min-w-0">
        <header className="md:hidden h-14 bg-white border-b border-slate-200 flex items-center px-4">
          <button onClick={() => setOpen(true)}>
            <Menu className="w-6 h-6" />
          </button>
        </header>
        <main className="p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
