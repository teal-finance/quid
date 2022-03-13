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
    require('@snowind/plugin')
  ],
  theme: {
    extend: {
      colors: {
        'primary': {
          DEFAULT: colors.cyan[700],
          dark: colors.cyan[900],
        },
      }
    },
  }
}