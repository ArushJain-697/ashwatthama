import { LABEL_CONFIRMED, LABEL_DRIFT, LABEL_HEURISTIC, LABEL_UNVERIFIABLE } from '../types'
import type { Zone } from '../graph'

const ITEMS: { zones: Zone[]; color: string; text: string }[] = [
  { zones: ['core', 'blast'], color: '#3fb950', text: `Core / Blast Radius — ${LABEL_CONFIRMED} (distance = real hop count)` },
  { zones: ['uncertain'], color: '#d29922', text: `Uncertain ring — ${LABEL_HEURISTIC}` },
  { zones: ['fog'], color: '#bc8cff', text: `Fog zone — ${LABEL_UNVERIFIABLE}` },
  { zones: ['outside'], color: '#f85149', text: `Outside the blast radius — ${LABEL_DRIFT}` },
]

export function Legend({
  visibleZones,
  onToggle,
}: {
  visibleZones: Set<string>
  onToggle: (zones: Zone[], checked: boolean) => void
}) {
  return (
    <div className="flex gap-4 mt-3.5 text-xs text-muted flex-wrap items-center">
      {ITEMS.map((item) => (
        <label key={item.text} className="inline-flex items-center gap-1.5 cursor-pointer select-none">
          <input
            type="checkbox"
            checked={item.zones.every((zone) => visibleZones.has(zone))}
            onChange={(event) => onToggle(item.zones, event.target.checked)}
            style={{ accentColor: '#ff4500' }}
          />
          <i
            className="w-2.5 h-2.5 rounded-full inline-block"
            style={{ background: item.color, boxShadow: `0 0 8px ${item.color}` }}
          />
          {item.text}
        </label>
      ))}
    </div>
  )
}
