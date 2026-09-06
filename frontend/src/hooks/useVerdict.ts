import { useEffect, useState } from 'react'
import type { Verdict } from '../types'

export type VerdictState =
  | { status: 'loading' }
  | { status: 'error'; message: string }
  | { status: 'ready'; verdict: Verdict }

// Fetches verdict.json at runtime from the same directory the app is served
// from -- in dev, Vite serves frontend/public/verdict.json (a sample fixture
// checked into the repo); in a real deployment, drop the real verdict.json
// from `verify-intent --out` next to the built dist/ bundle before serving
// it. Deliberately NOT bundled at build time: a rebuild-per-verdict workflow
// would defeat the point of a report you regenerate on every run.
export function useVerdict(path = './verdict.json'): VerdictState {
  const [state, setState] = useState<VerdictState>({ status: 'loading' })

  useEffect(() => {
    let cancelled = false
    fetch(path)
      .then((response) => {
        if (!response.ok) throw new Error(`${response.status} ${response.statusText}`)
        return response.json() as Promise<Verdict>
      })
      .then((verdict) => {
        if (!cancelled) setState({ status: 'ready', verdict })
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setState({
            status: 'error',
            message: error instanceof Error ? error.message : String(error),
          })
        }
      })
    return () => {
      cancelled = true
    }
  }, [path])

  return state
}
