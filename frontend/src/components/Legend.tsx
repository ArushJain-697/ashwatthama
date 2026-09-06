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
    <div className="gv-legend">
      {ITEMS.map((item) => (
        <label key={item.text}>
          <input
            type="checkbox"
            checked={item.zones.every((zone) => visibleZones.has(zone))}
            onChange={(event) => onToggle(item.zones, event.target.checked)}
          />
          <i style={{ background: item.color }} />
          {item.text}
        </label>
      ))}
    </div>
  )
}
