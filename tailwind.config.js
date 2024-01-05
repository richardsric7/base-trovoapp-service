/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
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
      },
    },
    fontFamily: {
      matahariRegular: ['MatahariRegular', 'sans-serif'],
      matahariExtended: ['MatahariExtended', 'sans-serif'],
      montserratMedium: ['MontserratMedium', 'sans-serif'],
      montserratRegular: ['MontserratRegular', 'sans-serif'],
      montserratSemiBold: ['MontserratSemiBold', 'sans-serif'],
    },
  },
  plugins: [],
};
