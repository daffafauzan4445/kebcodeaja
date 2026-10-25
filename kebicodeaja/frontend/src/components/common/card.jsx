export default function Card({ title, action, children, className = '' }) {
  return (
    <section className={`bg-white rounded-xl border border-slate-200 ${className}`}>
      {title && (
        <div className="flex items-center justify-between px-5 py-4 border-b border-slate-100">
          <h2 className="font-semibold text-slate-800">{title}</h2>
          {action}
        </div>
      )}
      <div className="p-5">{children}</div>
    </section>
  )
}
