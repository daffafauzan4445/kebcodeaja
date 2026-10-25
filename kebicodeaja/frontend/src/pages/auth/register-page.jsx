import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { User, Mail, Lock, UserPlus, Loader2, GraduationCap, Presentation } from 'lucide-react'
import { useAuth } from '../../context/auth-context'

const fields = [
  { name: 'name', label: 'Nama Lengkap', type: 'text', icon: User, placeholder: 'Nama lengkap' },
  { name: 'email', label: 'Email', type: 'email', icon: Mail, placeholder: 'nama@email.com' },
  { name: 'password', label: 'Password', type: 'password', icon: Lock, placeholder: 'Minimal 8 karakter' },
  { name: 'confirm', label: 'Konfirmasi Password', type: 'password', icon: Lock, placeholder: 'Ulangi password' },
]

const roles = [
  { value: 'student', label: 'Siswa', icon: GraduationCap },
  { value: 'teacher', label: 'Guru', icon: Presentation },
]

export default function RegisterPage() {
  const { register } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({ name: '', email: '', password: '', confirm: '', role: 'student' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    if (form.password.length < 8) return setError('Password minimal 8 karakter')
    if (form.password !== form.confirm) return setError('Konfirmasi password tidak sama')

    setLoading(true)
    try {
      await register({
        name: form.name,
        email: form.email,
        password: form.password,
        role: form.role,
      })
      navigate('/login', { replace: true, state: { registered: true } })
    } catch (err) {
      setError(err.response?.data?.message || 'Tidak dapat terhubung ke server')
    } finally {
      setLoading(false)
    }
  }

  const inputCls =
    'w-full pl-10 pr-3 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none'

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 to-slate-100 p-4">
      <div className="w-full max-w-md bg-white rounded-2xl shadow-lg p-8">
        <h1 className="text-2xl font-bold text-slate-800">Buat Akun</h1>
        <p className="text-sm text-slate-500 mt-1 mb-6">Bergabung dengan keBCodeAja</p>

        {error && (
          <div className="mb-4 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            {roles.map(({ value, label, icon: Icon }) => (
              <button
                type="button" key={value}
                onClick={() => setForm({ ...form, role: value })}
                className={`flex flex-col items-center gap-1 py-3 rounded-lg border-2 transition ${
                  form.role === value
                    ? 'border-indigo-600 bg-indigo-50 text-indigo-700'
                    : 'border-slate-200 text-slate-500 hover:border-slate-300'
                }`}
              >
                <Icon className="w-5 h-5" />
                <span className="text-sm font-medium">{label}</span>
              </button>
            ))}
          </div>

          {fields.map(({ name, label, type, icon: Icon, placeholder }) => (
            <div key={name}>
              <label className="block text-sm font-medium text-slate-700 mb-1">{label}</label>
              <div className="relative">
                <Icon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
                <input
                  type={type} name={name} required value={form[name]}
                  onChange={handleChange} placeholder={placeholder} className={inputCls}
                />
              </div>
            </div>
          ))}

          <button
            type="submit" disabled={loading}
            className="w-full flex items-center justify-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-60 text-white font-medium py-2.5 rounded-lg transition"
          >
            {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <UserPlus className="w-4 h-4" />}
            {loading ? 'Memproses...' : 'Daftar'}
          </button>
        </form>

        <p className="text-sm text-center text-slate-500 mt-6">
          Sudah punya akun?{' '}
          <Link to="/login" className="text-indigo-600 font-medium hover:underline">Masuk</Link>
        </p>
      </div>
    </div>
  )
}
