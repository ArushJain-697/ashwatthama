import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Relative base so the built dist/ can be dropped anywhere (a fidelity-out/
// subdirectory, a plain static file server, a CDN path) without rebuilding.
export default defineConfig({
  plugins: [react()],
  base: './',
})
