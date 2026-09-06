import { useEffect, useState } from 'react'

export function Header({ checkpointId, verdictLabel }: { checkpointId: string; verdictLabel: string }) {
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    function onScroll() {
      setScrolled(window.scrollY > 20)
    }
    window.addEventListener('scroll', onScroll)
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  const isReview = verdictLabel === 'REVIEW_REQUIRED'

  return (
    <header
      className={`sticky top-0 z-10 transition-all duration-400 border-b ${
        scrolled ? 'py-3 bg-black/85 backdrop-blur-md border-line' : 'py-5.5 border-transparent'
      }`}
    >
      <div className="px-7">
        <h1 className="font-serif text-[22px] font-semibold tracking-tight m-0">
          Fidelity<span className="text-accent">.</span>
        </h1>
        <div className="flex items-center gap-2.5 mt-1 text-xs text-muted font-mono">
          <span>{checkpointId}</span>
          <span
            className={`inline-flex items-center px-3 py-0.5 rounded-full text-[11px] uppercase tracking-wide border ${
              isReview ? 'text-[#f85149] border-[#f8514955] bg-[#f8514912]' : 'text-[#3fb950] border-[#3fb95055] bg-[#3fb95012]'
            }`}
          >
            {verdictLabel}
          </span>
        </div>
      </div>
    </header>
  )
}
