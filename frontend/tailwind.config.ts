import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        bg: '#0d1117',
        surface: '#161b22',
        border: '#30363d',
        text: {
          primary: '#e6edf3',
          muted: '#8b949e',
        },
        level: {
          INFO: '#3fb950',
          WARN: '#d29922',
          ERROR: '#f85149',
          DEBUG: '#58a6ff',
        },
        service: {
          auth: '#a371f7',
          order: '#f0883e',
          payment: '#39d353',
        },
      },
      fontFamily: {
        sans: ['Inter', 'Segoe UI', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Consolas', 'monospace'],
      },
      transitionDuration: {
        DEFAULT: '150ms',
      },
      transitionTimingFunction: {
        DEFAULT: 'ease',
      },
      boxShadow: {
        glow: '0 0 0 1px rgba(88, 166, 255, 0.25), 0 0 24px rgba(88, 166, 255, 0.15)',
      },
      keyframes: {
        pulseDot: {
          '0%, 100%': { opacity: '1', transform: 'scale(1)' },
          '50%': { opacity: '0.4', transform: 'scale(0.9)' },
        },
      },
      animation: {
        'pulse-dot': 'pulseDot 1.4s ease-in-out infinite',
      },
    },
  },
  plugins: [],
} satisfies Config
