import { useCallback, useEffect, useState } from 'react'
import { Users, GraduationCap, FileText, ClipboardList, Plus, Copy, Check, Loader2 } from 'lucide-react'
import api from '../../api/axios'
import errorMessage from '../../api/error-message'
import { useAuth } from '../../context/auth-context'
import PageHeader from '../../components/common/page-header'
import StatCard from '../../components/common/stat-card'
import ProgressBar from '../../components/common/progress-bar'
import Card from '../../components/common/card'
import Modal from '../../components/common/modal'
import { CT_ASPECTS } from '../../data/constants'

const emptyData = {
  stats: { classes: 0, students: 0, materials: 0, challenges: 0 },
  classes: [],
  averages: Object.keys(CT_ASPECTS).map((aspect) => ({ aspect, value: 0 })),
  recent_submissions: [],
}

const inputCls =
  'w-full px-3 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none'

export default function DashboardGuru() {
  const { user } = useAuth()
  const [data, setData] = useState(emptyData)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [form, setForm] = useState({ name: '', description: '' })
  const [saving, setSaving] = useState(false)
  const [formError, setFormError] = useState('')
  const [copied, setCopied] = useState('')

  const loadData = useCallback(async () => {
    try {
      const res = await api.get('/teacher/dashboard')
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

  const openModal = () => {
    setForm({ name: '', description: '' })
    setFormError('')
    setModalOpen(true)
  }

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleCreate = async (e) => {
    e.preventDefault()
    setFormError('')
    setSaving(true)
    try {
      const res = await api.post('/teacher/classes', form)
      const created = res.data.class
      setModalOpen(false)
      setNotice(`Kelas "${created.name}" berhasil dibuat. Kode kelas: ${created.class_code}`)
      await loadData()
    } catch (err) {
      setFormError(errorMessage(err, 'Gagal membuat kelas'))
    } finally {
      setSaving(false)
    }
  }

  const copyCode = async (code) => {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(code)
      setTimeout(() => setCopied(''), 1500)
    } catch {
      setCopied('')
    }
  }

  const today = new Date().toLocaleDateString('id-ID', {
    weekday: 'long', day: 'numeric', month: 'long', year: 'numeric',
  })

  const { stats, classes, averages, recent_submissions: submissions } = data

  if (loading) {
    return (
      <div className="flex items-center justify-center py-24 text-slate-400">
        <Loader2 className="w-6 h-6 animate-spin" />
      </div>
    )
  }

  return (
    <div>
      <PageHeader title={`Halo, ${user.name} `} subtitle={`${today} · Ringkasan kelas kamu`}>
        <button
          type="button"
          onClick={openModal}
          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-4 py-2.5 rounded-lg transition"
        >
          <Plus className="w-4 h-4" /> Buat Kelas
        </button>
      </PageHeader>

      {notice && (
        <div className="mb-4 rounded-lg bg-green-50 border border-green-200 text-green-700 text-sm px-4 py-3 flex items-start justify-between gap-3">
          <span>{notice}</span>
          <button type="button" onClick={() => setNotice('')} className="text-green-700 font-medium">
            Tutup
          </button>
        </div>
      )}

      {error && (
        <div className="mb-4 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
        <StatCard icon={Users} label="Total kelas" value={stats.classes} color="indigo" />
        <StatCard icon={GraduationCap} label="Total siswa" value={stats.students} color="emerald" />
        <StatCard icon={FileText} label="Materi" value={stats.materials} color="amber" />
        <StatCard icon={ClipboardList} label="Tantangan" value={stats.challenges} color="rose" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card title="Kelas Saya" className="lg:col-span-2">
          {classes.length === 0 ? (
            <p className="text-sm text-slate-500 text-center py-6">
              Belum ada kelas. Klik Buat Kelas untuk memulai.
            </p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {classes.map((c) => (
                <div key={c.id} className="border border-slate-200 rounded-lg p-4">
                  <p className="font-semibold text-slate-800">{c.name}</p>
                  <p className="text-sm text-slate-500 mt-1">
                    {c.student_count} siswa · {c.material_count} materi · {c.challenge_count} tantangan
                  </p>
                  <button
                    type="button"
                    onClick={() => copyCode(c.class_code)}
                    title="Klik untuk menyalin kode kelas"
                    className="mt-3 inline-flex items-center gap-2 bg-slate-100 hover:bg-slate-200 rounded-md px-2.5 py-1 text-xs font-mono text-slate-700 transition"
                  >
                    {copied === c.class_code ? (
                      <Check className="w-3 h-3 text-green-600" />
                    ) : (
                      <Copy className="w-3 h-3" />
                    )}
                    {copied === c.class_code ? 'Tersalin' : c.class_code}
                  </button>
                </div>
              ))}
            </div>
          )}
        </Card>

        <Card title="Rata-rata Nilai per Aspek CT">
          <div className="space-y-4">
            {averages.map((a) => (
              <ProgressBar
                key={a.aspect}
                label={CT_ASPECTS[a.aspect].label}
                value={a.value}
                barClass={CT_ASPECTS[a.aspect].bar}
              />
            ))}
          </div>
        </Card>

        <Card title="Pengumpulan Terbaru" className="lg:col-span-3">
          {submissions.length === 0 ? (
            <p className="text-sm text-slate-500 text-center py-6">Belum ada pengumpulan.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-slate-500 border-b border-slate-100">
                    <th className="pb-2 font-medium">Siswa</th>
                    <th className="pb-2 font-medium">Tantangan</th>
                    <th className="pb-2 font-medium">Kelas</th>
                    <th className="pb-2 font-medium">Nilai</th>
                    <th className="pb-2 font-medium">Waktu</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {submissions.map((s) => (
                    <tr key={s.id}>
                      <td className="py-3 font-medium text-slate-800">{s.student_name}</td>
                      <td className="py-3 text-slate-600">{s.challenge_title}</td>
                      <td className="py-3 text-slate-600">{s.class_name}</td>
                      <td className="py-3 font-semibold text-indigo-600">{s.score}</td>
                      <td className="py-3 text-slate-400">
                        {new Date(s.submitted_at).toLocaleString('id-ID', {
                          day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit',
                        })}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Card>
      </div>

      <Modal open={modalOpen} title="Buat Kelas Baru" onClose={() => setModalOpen(false)}>
        <form onSubmit={handleCreate} className="space-y-4">
          {formError && (
            <div className="rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm px-4 py-3">
              {formError}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Nama Kelas</label>
            <input
              type="text" name="name" required maxLength={150} value={form.name}
              onChange={handleChange} placeholder="Contoh: Informatika X-A" className={inputCls}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Deskripsi (opsional)</label>
            <textarea
              name="description" rows={3} maxLength={500} value={form.description}
              onChange={handleChange} placeholder="Keterangan singkat tentang kelas" className={inputCls}
            />
          </div>

          <p className="text-xs text-slate-500">Kode kelas dibuat otomatis setelah kelas disimpan.</p>

          <div className="flex justify-end gap-2">
            <button
              type="button" onClick={() => setModalOpen(false)}
              className="px-4 py-2.5 text-sm font-medium text-slate-600 hover:bg-slate-100 rounded-lg"
            >
              Batal
            </button>
            <button
              type="submit" disabled={saving}
              className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-60 text-white text-sm font-medium px-4 py-2.5 rounded-lg transition"
            >
              {saving && <Loader2 className="w-4 h-4 animate-spin" />}
              {saving ? 'Menyimpan...' : 'Simpan'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
