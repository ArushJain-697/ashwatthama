/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        bg: '#171e19',
        card: '#ffffff',
        line: '#000000',
        muted: '#5b645f',
        accent: '#ffe17c',
        sage: '#b7c6c2',
        charcoal: '#171e19',
        tier: {
          confirmed: '#3fb950',
          heuristic: '#d29922',
          unverifiable: '#bc8cff',
          drift: '#f85149',
        },
      },
      fontFamily: {
        sans: ['Satoshi', 'ui-sans-serif', 'system-ui'],
        serif: ['"Cabinet Grotesk"', 'ui-sans-serif'],
        mono: ['ui-monospace', 'monospace'],
      },
      boxShadow: {
        hard: '8px 8px 0px 0px #000000',
        'hard-sm': '4px 4px 0px 0px #000000',
        'hard-lg': '12px 12px 0px 0px #000000',
      },
      backgroundImage: {
        'dot-pattern':
          'radial-gradient(circle, rgba(0,0,0,0.1) 1.5px, transparent 1.5px)',
      },
      backgroundSize: {
        dots: '32px 32px',
      },
    },
  },
  plugins: [],
}
