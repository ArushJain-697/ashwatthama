import { useEffect, useRef, useState } from 'react'

// Reveal-on-scroll, with a fix for a real bug found in the Go-rendered
// single-file dashboard's version of this same pattern: there, elements
// started at opacity:0 and depended ENTIRELY on an IntersectionObserver
// callback to ever become visible. In at least one real browser session,
// that callback never fired, and the content -- including the four summary
// cards -- stayed permanently invisible ("a black screen").
//
// This hook keeps the same visual effect but never lets the observer be a
// single point of failure: a mount-time fallback timer marks the element
// active regardless, after a short delay, so a browser quirk degrades to
// "the fade-in didn't play" instead of "the content never appeared."
export function useReveal<T extends HTMLElement>(fallbackDelayMs = 800) {
  const ref = useRef<T | null>(null)
  const [active, setActive] = useState(false)

  useEffect(() => {
    const element = ref.current
    if (!element) return

    const fallback = window.setTimeout(() => setActive(true), fallbackDelayMs)

    if (typeof IntersectionObserver === 'undefined') {
      // No observer support at all -- show immediately rather than wait.
      setActive(true)
      return () => window.clearTimeout(fallback)
    }

    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            setActive(true)
            window.clearTimeout(fallback)
          }
        }
      },
      { threshold: 0.1, rootMargin: '0px 0px -40px 0px' },
    )
    observer.observe(element)
    return () => {
      window.clearTimeout(fallback)
      observer.disconnect()
    }
  }, [fallbackDelayMs])

  return { ref, active }
}
