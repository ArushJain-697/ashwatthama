# Fidelity dashboard (React)

A real component-based frontend for `verdict.json`, separate from the
single-file `dashboard.html` the Go CLI generates directly
(`internal/fidelity/dashboard.go`'s `RenderDashboard`). That file stays and
keeps working — it's what `verify-intent` writes on every run, and its
zero-dependency, double-click-to-open property is a real requirement (no
server, no build step, always in sync with the run that produced it). This
app is a richer, separate frontend for development or a more polished demo,
built with React + TypeScript + Vite + Tailwind + `@react-three/fiber`.

## What it is

Same data, same six-tier model, same radial blast-radius 3D scene as the
Go-rendered dashboard — reimplemented as real components instead of one
generated HTML blob, so it's actually maintainable as a UI:

- `src/types.ts` — TypeScript types mirroring `verdict.go`'s `Verdict` struct
  field-for-field (kept in sync by hand, no generator).
- `src/graph.ts` — a faithful port of `dashboard.go`'s `dashboardGraph`: real
  edges exist only for `expected_blast_radius`; every other tier is an
  unconnected, tier-colored node placed in its own radial zone.
- `src/components/graph/` — the Three.js scene as `@react-three/fiber`
  components (`Nodes`, `Edges`, `HopRings`, `BlastRadiusScene`), using drei's
  `OrbitControls`/`Stars`/`Sparkles` and `@react-three/postprocessing`'s
  `Bloom` instead of hand-rolled Three.js — the same visual result as the
  Go-rendered scene, written as idiomatic React instead of imperative
  `useEffect` DOM manipulation.
- `src/components/SummaryTab.tsx`, `TierCard.tsx`, `TierList.tsx` — the
  tier-grouped summary view.
- `src/hooks/useReveal.ts` — reveal-on-scroll, **with a fix for a real bug**
  found in the Go-rendered version's equivalent: there, elements started
  `opacity:0` and depended entirely on an `IntersectionObserver` callback to
  ever appear, and in one real browser session that callback never fired —
  the content stayed permanently invisible. This hook adds a mount-time
  fallback timer so a browser quirk degrades to "no fade-in animation"
  instead of "the content never appeared."

## Dev

```sh
npm install
npm run dev
```

Reads `public/verdict.json` (a real fixture, copied in from a `verify-intent`
run — replace it with a fresher one any time; the page hot-reloads).

## Build and deploy

```sh
npm run build   # -> dist/, ~1MB (Three.js is not small; not code-split, not asked for)
```

`dist/` is a static bundle with relative asset paths (`vite.config.ts`'s
`base: './'`), so it can be dropped anywhere and served from a plain static
file server. It does **not** work by double-clicking `dist/index.html`
directly — `fetch('./verdict.json')` is blocked by browsers under `file://`
CORS rules, the same reason the Go-rendered dashboard embeds its data inline
instead. Serve it:

```sh
cp path/to/real/verdict.json dist/verdict.json
cd dist && python3 -m http.server 8080   # or any static file server
```

Verified end-to-end this way (headless Chromium DOM dump, zero console
errors, all reveal-on-scroll elements activated) before this was committed.

## What's not done

- No code-splitting (`vite build`'s own size warning) — Three.js dominates
  the bundle; not worth the complexity for a report-viewing tool that loads
  once per session.
- No automated pipeline connecting `verify-intent`'s output to this app's
  `dist/` — copying `verdict.json` next to the built bundle is a manual step,
  documented above rather than scripted, since adding a hard Node/npm
  runtime dependency to the `entire-graph` Go binary itself would break its
  single-native-binary distribution model.
