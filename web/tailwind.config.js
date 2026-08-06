/** @type {import('tailwindcss').Config} */
const path = require('path');

module.exports = {
  content: [
    './index.html',
    './src/**/*.{js,jsx,ts,tsx}',
    './node_modules/tw-elements-react/dist/js/**/*.js',
    path.join(__dirname, './pages/**/*.{js,ts,jsx,tsx}'),
    path.join(__dirname, './components/**/*.{js,ts,jsx,tsx}'),
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          100: '#F2F6F9',
          200: '#CCDBE7',
          300: '#99B6CF',
          500: '#99B6CF',
          600: '#6692B8',
          700: '#336DA0',
          800: '#004988',
        },
        trovored: {
          light: '#FFE7E2',
          primary: '#BE3800',
        },
      },
    },
    fontFamily: {
      matahariRegular: ['MatahariRegular', 'sans-serif'],
      matahariExtended: ['MatahariExtended', 'sans-serif'],
      montserratMedium: ['MontserratMedium', 'sans-serif'],
      montserratRegular: ['MontserratRegular', 'sans-serif'],
      montserratSemiBold: ['MontserratSemiBold', 'sans-serif'],
    },
    fontWeight: {
      bold: '700', // Use '700' for Montserrat-Bold
    },
  },
  plugins: [require('tw-elements-react/dist/plugin.cjs')],
};
