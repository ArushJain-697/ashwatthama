import type { ReactNode } from 'react'
import { useReveal } from '../hooks/useReveal'

export function TierCard({
  count,
  label,
  color,
  delayMs = 0,
}: {
  count: number
  label: ReactNode
  color: string
  delayMs?: number
}) {
  const { ref, active } = useReveal<HTMLDivElement>()
  return (
    <div
      ref={ref}
      data-active={active}
      className="reveal bg-card border border-line rounded-[20px] p-5.5 transition-all duration-400 ease-[cubic-bezier(.22,1,.36,1)] hover:-translate-y-1"
      style={
        {
          transitionDelay: `${delayMs}ms`,
          '--tw-shadow': `0 16px 40px -12px ${color}40`,
        } as React.CSSProperties
      }
      onMouseEnter={(e) => {
        e.currentTarget.style.borderColor = color
        e.currentTarget.style.boxShadow = `0 16px 40px -12px ${color}40`
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.borderColor = ''
        e.currentTarget.style.boxShadow = ''
      }}
    >
      <div className="font-serif text-4xl font-semibold leading-none" style={{ color }}>
        {count}
      </div>
      <div className="text-xs text-muted mt-2.5 uppercase tracking-wide">{label}</div>
    </div>
  )
}
