/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        bg: '#14161A',
        surface: '#1C1F26',
        surface2: '#232730',
        border: '#2A2E37',
        text: '#E8E9ED',
        muted: '#8B8F99',
        accent: '#5B8DEF',
        'accent-dim': '#3D5FA8',
        get: '#5B8DEF',
        post: '#4ADE80',
        put: '#FBBF24',
        patch: '#C084FC',
        del: '#F87171',
        ok: '#4ADE80',
        redirect: '#60A5FA',
        clienterr: '#FBBF24',
        servererr: '#F87171',
      },
      fontFamily: {
        display: ['"Space Grotesk"', 'sans-serif'],
        sans: ['Inter', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      },
    },
  },
  plugins: [],
}
