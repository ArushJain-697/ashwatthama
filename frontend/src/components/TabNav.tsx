export type Tab = 'summary' | 'graph'

export function TabNav({ tab, onChange }: { tab: Tab; onChange: (tab: Tab) => void }) {
  const tabs: { id: Tab; label: string }[] = [
    { id: 'summary', label: 'Summary' },
    { id: 'graph', label: 'Graph View' },
  ]
  return (
    <nav className="flex gap-3 px-7 py-5 bg-charcoal">
      {tabs.map((t) => {
        const active = t.id === tab
        return (
          <button
            key={t.id}
            onClick={() => onChange(t.id)}
            className={`brutal-btn font-sans font-medium text-sm px-5 py-2 border-2 border-black ${
              active ? 'bg-black text-accent shadow-hard-sm' : 'bg-white text-black shadow-hard-sm'
            }`}
          >
            {t.label}
          </button>
        )
      })}
    </nav>
  )
}
