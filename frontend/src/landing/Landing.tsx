import { useEffect } from 'react'
import './landing.css'
import { navigate } from '../App'

// Ported from Dwarpal's Landing.jsx — same Neo-Brutalist sections, copy
// rewritten for Fidelity: intent-vs-diff verification reconciled by a code
// graph, six tiers, one verdict.json.

function Arw() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M4 12 12 4M6 4h6v6" />
    </svg>
  )
}

type BtnProps = { size?: 'lg' | 'sm'; variant?: 'black' | 'white' | 'yellow'; children?: React.ReactNode }

function OpenBtn({ size = 'lg', variant = 'black', children = 'Open the dashboard' }: BtnProps) {
  return (
    <button className={`btn btn-${variant} btn-${size}`} onClick={() => navigate('/dashboard')}>
      {children}
      <span className="arw">
        <Arw />
      </span>
    </button>
  )
}

function AnchorBtn({ href, size = 'sm', variant = 'white', children }: BtnProps & { href: string }) {
  return (
    <a className={`btn btn-${variant} btn-${size}`} href={href}>
      {children}
      <span className="arw">
        <Arw />
      </span>
    </a>
  )
}

function Nav() {
  return (
    <nav className="bnav">
      <div className="bnav-in">
        <a className="brand" href="/" onClick={(e) => { e.preventDefault(); navigate('/') }}>
          <span className="brand-mark">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ffe17c" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M12 3 4 6v5c0 5 3.5 8 8 10 4.5-2 8-5 8-10V6l-8-3Z" />
              <path d="m9 11 2 2 4-4" />
            </svg>
          </span>
          Fidelity
        </a>
        <div className="bnav-links">
          <a href="#how">How it works</a>
          <a href="#tiers">Tiers</a>
          <a href="#verdict">Verdict</a>
          <a href="/dashboard" onClick={(e) => { e.preventDefault(); navigate('/dashboard') }}>Dashboard</a>
        </div>
        <button className="btn btn-black btn-sm" onClick={() => navigate('/dashboard')}>
          Open a verdict →
        </button>
      </div>
    </nav>
  )
}

function Hero() {
  return (
    <header className="hero dots">
      <div className="wrap hero-in">
        <div>
          <span className="hbadge">
            <span className="sq" />
            <span>Intent · diff · verdict</span>
          </span>
          <h1>
            Did it do <span className="stroke">only</span>
            <br />
            what you asked?
          </h1>
          <p>
            Fidelity reconciles the prompt behind a change against the diff that landed — over a real
            code graph — and returns a verdict you can gate on.
          </p>
          <div className="hero-cta">
            <OpenBtn size="lg" variant="black" />
            <AnchorBtn href="#how" size="sm" variant="white">
              How it works
            </AnchorBtn>
          </div>
          <p className="hero-hint">Reads the verdict.json from your last verify-intent run — no setup.</p>
        </div>

        <div className="mock">
          <div className="sticker">
            no silent
            <br />
            scope creep
          </div>
          <div className="mock-bar">
            <i className="r" />
            <i className="y" />
            <i className="g" />
            <span>fidelity · verdict</span>
          </div>
          <div className="mock-body">
            <div className="mock-chips">
              <div className="mock-chip yell">
                <small>Confirmed</small>
                <b>65</b>
              </div>
              <div className="mock-chip dark">
                <small>Scope creep</small>
                <b>620</b>
              </div>
              <div className="mock-chip sage">
                <small>Advisory</small>
                <b>19</b>
              </div>
            </div>
            {[
              ['CONFIRMED', 'BuildVerdict', '✓', '#1a9d5c'],
              ['BLAST', 'ComputeCoverage', '2 hops', '#171e19'],
              ['ADVISORY', 'Corroboration', '?', '#b5820f'],
              ['CREEP', 'DefaultConfig', '!', '#d21f3c'],
            ].map(([cat, entity, mark, color]) => (
              <div className="mock-row" key={cat}>
                <span className="cat">{cat}</span>
                <span>{entity}</span>
                <span className="code" style={{ color }}>
                  {mark}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </header>
  )
}

const MARQUEE = [
  'CONFIRMED',
  'EXPECTED BLAST RADIUS',
  'ADVISORY',
  'UNVERIFIABLE',
  'SCOPE CREEP',
  'UNIMPLEMENTED',
  'CALLS',
  'EXTENDS',
  'USES_TYPE',
]

function Marquee() {
  const items = [...MARQUEE, ...MARQUEE]
  return (
    <div className="marquee">
      <div className="marquee-track">
        {items.map((m, i) => (
          <span key={i}>{m} ✳</span>
        ))}
      </div>
    </div>
  )
}

const PROBLEM = [
  'Skim the diff, trust the PR description',
  'No check that the stated intent was actually implemented',
  'Silent scope creep slips past review',
  'Blast radius guessed from memory',
  '“Looks fine” — with no record of why',
]
const SOLUTION = [
  'Extracts the intent, matches it to changed entities',
  'Flags what was declared but never implemented',
  'Names every change with no path back to the ask',
  'Blast radius is a real hop count on the code graph',
  'Every run writes a verdict.json you can gate CI on',
]

function ProblemSolution() {
  return (
    <section className="sec sec-pad ps">
      <div className="wrap">
        <h2 className="sec-title">Reviewing a diff by eye misses the two things that matter.</h2>
        <p className="sec-sub">Whether it did what was asked — and whether it did more.</p>
        <div className="ps-grid">
          <div className="ps-card problem">
            <h3>The usual way</h3>
            <ul className="ps-list">
              {PROBLEM.map((t) => (
                <li key={t}>
                  <span className="ps-ic">✕</span>
                  {t}
                </li>
              ))}
            </ul>
          </div>
          <div className="ps-card solution">
            <h3>The Fidelity way</h3>
            <ul className="ps-list">
              {SOLUTION.map((t) => (
                <li key={t}>
                  <span className="ps-ic">✓</span>
                  {t}
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  )
}

const FEATURES = [
  {
    h: 'Intent extraction',
    p: 'Pulls the named entities out of the prompt behind the change, each with the method that found it and a confidence.',
  },
  {
    h: 'Graph-reconciled blast radius',
    p: 'Walks real CALLS / EXTENDS / USES_TYPE edges from each confirmed change — distance is a measured hop count, not a guess.',
  },
  {
    h: 'Confidence gating',
    p: 'Low-confidence matches pass a conformal gate or abstain. The report never launders a guess as a fact.',
  },
  {
    h: 'Coverage honesty',
    p: "Where the graph couldn't see clearly, the entity lands in 'unverifiable' — with the source or test that would settle it.",
  },
  {
    h: 'Scope-creep detection',
    p: 'Every changed entity with no path back to the stated intent is surfaced, with its community label.',
  },
  {
    h: 'A verdict, not a vibe',
    p: 'One verdict.json per run: six tiers, a label, and the checkpoint it was computed against. Diff it, gate on it, archive it.',
  },
]

function Features() {
  return (
    <section className="sec sec-pad features dots" id="tiers">
      <div className="wrap">
        <h2 className="sec-title">Built to be checkable, not trusted.</h2>
        <p className="sec-sub">Every claim in the verdict traces to a structural fact or is labelled as not one.</p>
        <div className="feat-grid">
          {FEATURES.map((f) => (
            <div className="feat" key={f.h}>
              <h3>{f.h}</h3>
              <p>{f.p}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

const STEPS = [
  { n: '01', h: 'Extract', p: 'Read the prompt behind the change; resolve it to named entities in the repo.' },
  { n: '02', h: 'Reconcile', p: 'Diff the checkpoints, then reconcile every changed entity against the intent over the code graph.' },
  { n: '03', h: 'Judge', p: 'Sort each entity into a tier, gate the shaky ones, and emit a verdict with a label.' },
]

function HowItWorks() {
  return (
    <section className="sec sec-pad how" id="how">
      <div className="wrap">
        <h2 className="sec-title">One change, three moves.</h2>
        <p className="sec-sub">From the prompt that asked for it to a verdict you can act on.</p>
        <div className="steps">
          {STEPS.map((s) => (
            <div className="step" key={s.n}>
              <div className="step-num">{s.n}</div>
              <h3>{s.h}</h3>
              <p>{s.p}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

function Personas() {
  return (
    <section className="sec sec-pad who">
      <div className="wrap">
        <h2 className="sec-title">For the moment before you merge.</h2>
        <p className="sec-sub">Review, ship, and audit off the same verdict.</p>
        <div className="who-grid">
          <div className="persona sage">
            <span className="pill">Review</span>
            <h3>See it in one screen</h3>
            <p>Whether a change matches its ask — and what it touched that nobody mentioned.</p>
          </div>
          <div className="persona yell">
            <span className="pill">Ship</span>
            <h3>Gate CI on the label</h3>
            <p>
              <code>REVIEW_REQUIRED</code> stops the merge until a human has looked at the drift.
            </p>
          </div>
          <div className="persona dark">
            <span className="pill">Audit</span>
            <h3>A record, not a memory</h3>
            <p>Every verdict.json is dated and checkpoint-anchored: what a change claimed, and what it did.</p>
          </div>
        </div>
      </div>
    </section>
  )
}

const DEMOS = [
  { h: 'Graph', items: ['real edge traversal', 'hop-count blast radius', 'community detection'] },
  { h: 'Calibration', items: ['conformal abstention', 'verbal confidence', 'corroboration checks'] },
  { h: 'Reporting', items: ['six-tier reconciliation', 'stable verdict schema', '3D blast-radius view'] },
]

function Demonstrates() {
  return (
    <section className="sec sec-pad demo-sec">
      <div className="wrap">
        <h2 className="sec-title">What it leans on.</h2>
        <p className="sec-sub">The pieces under the verdict, grouped by the job they do.</p>
        <div className="demo-grid">
          {DEMOS.map((d) => (
            <div className="demo-card" key={d.h}>
              <h3>{d.h}</h3>
              <ul>
                {d.items.map((i) => (
                  <li key={i}>
                    <span className="dot2" />
                    {i}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

function Command() {
  return (
    <section className="sec sec-pad cmd-sec" id="verdict">
      <div className="wrap">
        <h2 className="sec-title">See it in one command.</h2>
        <p className="sec-sub">Point verify-intent at two checkpoints and open the report it writes.</p>
        <div className="cmd">
          <div className="cmd-bar">
            <i style={{ background: '#ff5f57' }} />
            <i style={{ background: '#febc2e' }} />
            <i style={{ background: '#28c840' }} />
          </div>
          <div className="cmd-body">
            <div>
              <span className="p">$ </span>verify-intent --base HEAD~1 --head HEAD --out fidelity-out/
            </div>
            <div className="c">{'  '}extracting intent … 7 entities</div>
            <div className="c">{'  '}reconciling 158 changed entities over the graph</div>
            <div className="c">
              {'  '}confirmed 65 · advisory 19 · unverifiable 73 · scope-creep 620
            </div>
            <div className="c">
              {'  '}verdict: <span className="blk">REVIEW_REQUIRED</span> → fidelity-out/verdict.json
            </div>
            <div>
              <span className="p">$ </span>open fidelity-out/index.html<span className="c"> # the dashboard, offline</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

function FinalCta() {
  return (
    <section className="sec sec-pad final dots">
      <div className="wrap">
        <h2>Open the last verdict.</h2>
        <p>Six tiers, one label, and the graph it was computed on.</p>
        <div className="hero-cta">
          <OpenBtn size="lg" variant="black" />
          <AnchorBtn href="#how" size="sm" variant="white">
            How it works
          </AnchorBtn>
        </div>
      </div>
    </section>
  )
}

function Footer() {
  return (
    <footer className="foot2">
      <div className="wrap">
        <div className="foot2-grid">
          <div>
            <h4>Fidelity</h4>
            <p style={{ maxWidth: '32ch', lineHeight: 1.5 }}>
              Intent-vs-diff verification for code changes, reconciled by a local code graph and
              distilled to one verdict.json.
            </p>
          </div>
          <div className="cols">
            <h5>Product</h5>
            <a href="#how">How it works</a>
            <a href="#tiers">Tiers</a>
            <a href="/dashboard" onClick={(e) => { e.preventDefault(); navigate('/dashboard') }}>
              Dashboard
            </a>
          </div>
          <div className="cols">
            <h5>Built on</h5>
            <span>Go · verify-intent</span>
            <span>entire-graph</span>
            <span>React · Three.js</span>
            <span>verdict.json</span>
          </div>
        </div>
        <div className="foot2-bottom">Fidelity — did it do only what you asked?</div>
      </div>
    </footer>
  )
}

export function Landing() {
  useEffect(() => {
    const prev = document.body.style.background
    document.body.style.background = '#ffe17c'
    return () => {
      document.body.style.background = prev
    }
  }, [])

  return (
    <div className="brut">
      <Nav />
      <Hero />
      <Marquee />
      <ProblemSolution />
      <Features />
      <HowItWorks />
      <Personas />
      <Demonstrates />
      <Command />
      <FinalCta />
      <Footer />
    </div>
  )
}
