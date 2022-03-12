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
  /*theme: {
    extend: {
      colors: {
        'background': {
          DEFAULT: colors.white,
          dark: colors.neutral[900]
        },
        'primary': {
          DEFAULT: colors.cyan[700],
          dark: colors.cyan[800],
        },
        'secondary': {
          DEFAULT: colors.cyan[500],
          dark: colors.slate[600],
        },
        'light': {
          DEFAULT: colors.slate[200],
          dark: colors.neutral[700]
        },
      }
    },
  }*/
}