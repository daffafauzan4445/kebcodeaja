import { Construction } from 'lucide-react'
import PageHeader from './page-header'

export default function ComingSoon({ title }) {
  return (
    <div>
      <PageHeader title={title} />
      <div className="bg-white rounded-xl border border-dashed border-slate-300 p-12 text-center">
        <Construction className="w-10 h-10 text-slate-400 mx-auto mb-3" />
        <p className="font-medium text-slate-700">Halaman ini sedang dikembangkan</p>
        <p className="text-sm text-slate-500 mt-1">Fitur akan tersedia di pertemuan berikutnya.</p>
      </div>
    </div>
  )
}
