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
            bg: '#114B5E',
            txt: colors.neutral[100]
          }
        },
        topbar: {
          light: {
            bg: colors.cyan[800],
            txt: colors.white
          },
          dark: {
            bg: '#080807',
            txt: colors.neutral[400]
          }
        },
        sidebar: {
          light: {
            bg: colors.cyan[500],
            txt: colors.white
          },
          dark: {
            bg: '#171814',
            txt: colors.neutral[400]
          }
        },
      }
    },
  }
}