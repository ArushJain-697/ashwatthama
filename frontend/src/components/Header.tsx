export function Header({ checkpointId, verdictLabel }: { checkpointId: string; verdictLabel: string }) {
  const isReview = verdictLabel === 'REVIEW_REQUIRED'

  return (
    <header className="sticky top-0 z-10 h-20 bg-accent border-b-2 border-black bg-dot-pattern bg-dots flex items-center">
      <div className="px-7 flex items-center gap-4 w-full">
        <div className="w-10 h-10 bg-black border-2 border-black flex items-center justify-center shrink-0">
          <span className="text-accent font-serif font-extrabold text-lg leading-none">F</span>
        </div>
        <div>
          <h1 className="font-serif font-extrabold text-xl tracking-tight text-black m-0 leading-none">Fidelity</h1>
          <div className="flex items-center gap-2 mt-1 text-xs font-sans font-medium text-black/70">
            <span className="font-mono">{checkpointId}</span>
            <span
              className={`inline-flex items-center px-3 py-0.5 rounded-full text-[11px] uppercase tracking-wide border-2 border-black bg-white font-sans font-medium ${
                isReview ? 'text-[#c62828]' : 'text-[#1b7a34]'
              }`}
            >
              {verdictLabel}
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
