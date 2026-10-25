import { useCallback, useEffect, useState } from 'react'
import { Users, BookOpen, Trophy, Award, Clock, FileText, LogIn, Loader2 } from 'lucide-react'
import api from '../../api/axios'
import errorMessage from '../../api/error-message'
import { useAuth } from '../../context/auth-context'
import PageHeader from '../../components/common/page-header'
import StatCard from '../../components/common/stat-card'
import ProgressBar from '../../components/common/progress-bar'
import Card from '../../components/common/card'
import { CT_ASPECTS, DIFFICULTY } from '../../data/constants'

const emptyData = {
  stats: { classes: 0, materials: 0, completed_challenges: 0, average_score: 0 },
  progress: Object.keys(CT_ASPECTS).map((aspect) => ({ aspect, value: 0 })),
  challenges: [],
  materials: [],
}

export default function DashboardSiswa() {
  const { user } = useAuth()
  const [data, setData] = useState(emptyData)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [code, setCode] = useState('')
  const [joining, setJoining] = useState(false)
  const [joinError, setJoinError] = useState('')
  const [notice, setNotice] = useState('')

  const loadData = useCallback(async () => {
    try {
      const res = await api.get('/student/dashboard')
      setData(res.data)
      setError('')
    } catch (err) {
      setError(errorMessage(err, 'Gagal memuat data dashboard'))
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [loadData])

  const handleJoin = async (e) => {
    e.preventDefault()
    setJoinError('')
    setNotice('')
    setJoining(true)
    try {
      const res = await api.post('/student/classes/join', { class_code: code })
      setNotice(res.data.message)
      setCode('')
      await loadData()
    } catch (err) {
      setJoinError(errorMessage(err, 'Gagal bergabung ke kelas'))
    } finally {
      setJoining(false)
    }
  }

  const today = new Date().toLocaleDateString('id-ID', {
    weekday: 'long', day: 'numeric', month: 'long', year: 'numeric',
  })

  const { stats, progress, challenges, materials } = data

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24 text-slate-400">
        <Loader2 className="w-6 h-6 animate-spin" />
      </div>
    )
  }

  return (
    <div>
      <PageHeader title={`Halo, ${user.name} `} subtitle={`${today} · Semangat belajar hari ini!`} />

      {error && (
        <div className="mb-4 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3">
          {error}
        </div>
      )}

      <Card title="Gabung Kelas" className="mb-6">
        <form onSubmit={handleJoin} className="flex flex-col sm:flex-row gap-3">
          <input
            type="text" required value={code} maxLength={12}
            onChange={(e) => setCode(e.target.value.toUpperCase())}
            placeholder="Masukkan kode kelas, contoh: KBA-7H2Q"
            className="flex-1 px-3 py-2.5 border border-slate-300 rounded-lg font-mono focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
          />
          <button
            type="submit" disabled={joining}
            className="inline-flex items-center justify-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-60 text-white text-sm font-medium px-5 py-2.5 rounded-lg transition"
          >
            {joining ? <Loader2 className="w-4 h-4 animate-spin" /> : <LogIn className="w-4 h-4" />}
            {joining ? 'Memproses...' : 'Gabung'}
          </button>
        </form>

        {joinError && <p className="mt-3 text-sm text-red-600">{joinError}</p>}
        {notice && <p className="mt-3 text-sm text-green-600">{notice}</p>}
        <p className="mt-3 text-xs text-slate-500">Minta kode kelas kepada gurumu.</p>
      </Card>

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
        <StatCard icon={Users} label="Kelas diikuti" value={stats.classes} color="indigo" />
        <StatCard icon={BookOpen} label="Materi tersedia" value={stats.materials} color="emerald" />
        <StatCard icon={Trophy} label="Tantangan selesai" value={stats.completed_challenges} color="amber" />
        <StatCard icon={Award} label="Rata-rata nilai" value={stats.average_score} color="rose" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card title="Progres Berpikir Komputasional" className="lg:col-span-1">
          <div className="space-y-4">
            {progress.map((p) => (
              <ProgressBar
                key={p.aspect}
                label={CT_ASPECTS[p.aspect].label}
                value={p.value}
                barClass={CT_ASPECTS[p.aspect].bar}
              />
            ))}
          </div>
        </Card>

        <Card title="Tantangan Terbaru" className="lg:col-span-2">
          {challenges.length === 0 ? (
            <p className="text-sm text-slate-500 text-center py-6">Belum ada tantangan.</p>
          ) : (
            <ul className="divide-y divide-slate-100">
              {challenges.map((c) => (
                <li key={c.id} className="py-3 first:pt-0 last:pb-0 flex items-center justify-between gap-3">
                  <div className="min-w-0">
                    <p className="font-medium text-slate-800 truncate">{c.title}</p>
                    <div className="flex flex-wrap items-center gap-2 mt-1">
                      <span className="text-xs text-slate-500">{CT_ASPECTS[c.ct_aspect]?.label}</span>
                      <span className={`text-xs px-2 py-0.5 rounded-full ${DIFFICULTY[c.difficulty]?.cls}`}>
                        {DIFFICULTY[c.difficulty]?.label}
                      </span>
                      <span className="text-xs text-slate-400">{c.class_name}</span>
                    </div>
                  </div>
                  <div className="text-right shrink-0">
                    {c.score !== null ? (
                      <span className="text-sm font-semibold text-green-600">Nilai {c.score}</span>
                    ) : (
                      <span className="text-xs text-slate-400">Belum dikerjakan</span>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </Card>

        <Card title="Materi Terbaru" className="lg:col-span-3">
          {materials.length === 0 ? (
            <p className="text-sm text-slate-500 text-center py-6">Belum ada materi.</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {materials.map((m) => (
                <div key={m.id} className="border border-slate-200 rounded-lg p-4 hover:border-indigo-300 transition">
                  <FileText className="w-5 h-5 text-indigo-500 mb-2" />
                  <p className="font-medium text-slate-800">{m.title}</p>
                  <p className="text-sm text-slate-500 mt-1">{m.class_name}</p>
                  <p className="text-xs text-slate-400 mt-2 flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    {new Date(m.created_at).toLocaleDateString('id-ID', {
                      day: 'numeric', month: 'short', year: 'numeric',
                    })}
                  </p>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>
    </div>
  )
}
