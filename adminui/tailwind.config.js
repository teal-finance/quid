const colors = require('tailwindcss/colors')

module.exports = {
  content: [
    './index.html',
    './src/**/*.{js,jsx,ts,tsx,vue}',
    './node_modules/@snowind/**/*.{vue,js,ts}',
  ],
  darkMode: 'class',
  plugins: [
    require('@tailwindcss/forms'),
    require('@snowind/plugin'),
    require('tailwindcss-semantic-colors'),
  ],
  theme: {
    extend: {
      semanticColors: {
        primary: {
          light: {
            bg: colors.cyan[800],
            txt: colors.white
          },
          dark: {
            bg: '#0D3846',
            txt: colors.neutral[100]
          }
        }
      }
    },
  }
}