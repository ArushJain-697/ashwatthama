export type Tab = 'summary' | 'graph'

export function TabNav({ tab, onChange }: { tab: Tab; onChange: (tab: Tab) => void }) {
  const tabs: { id: Tab; label: string }[] = [
    { id: 'summary', label: 'Summary' },
    { id: 'graph', label: 'Graph View' },
  ]
  return (
    <nav className="flex gap-1.5 px-7 pb-4 relative z-10">
      {tabs.map((t) => {
        const active = t.id === tab
        return (
          <button
            key={t.id}
            onClick={() => onChange(t.id)}
            className={`border rounded-full px-5 py-2 text-[13px] transition-all duration-200 ${
              active
                ? 'text-white bg-accent border-accent shadow-[0_0_20px_rgba(255,69,0,0.35)]'
                : 'text-muted border-line hover:border-[#444] hover:text-fg bg-transparent'
            }`}
          >
            {t.label}
          </button>
        )
      })}
    </nav>
  )
}
