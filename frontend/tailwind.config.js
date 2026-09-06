/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        bg: '#050505',
        card: '#0e0e0e',
        line: '#232323',
        muted: '#8a8a8a',
        accent: '#ff4500',
        tier: {
          confirmed: '#3fb950',
          heuristic: '#d29922',
          unverifiable: '#bc8cff',
          drift: '#f85149',
        },
      },
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui'],
        serif: ['"Playfair Display"', 'serif'],
        mono: ['ui-monospace', 'monospace'],
      },
    },
  },
  plugins: [],
}
