const js = require('@eslint/js');
const react = require('eslint-plugin-react');
const reactHooks = require('eslint-plugin-react-hooks');
const prettier = require('eslint-config-prettier');
const globals = require('globals');

// typescript-eslint doesn't support TypeScript 7 yet (confirmed by a hard
// runtime guard in @typescript-eslint/parser itself, not just a peerDep
// warning - see https://github.com/typescript-eslint/typescript-eslint/issues/10940).
// Parsing here goes through @babel/eslint-parser + @babel/preset-typescript
// instead: syntax-only (strips types, no type-aware rules), independent of
// the `typescript` package version, so tsc stays on TypeScript 7 for the
// real type-checking gate (`npm run build`) while linting still works.
module.exports = [
  {
    ignores: ['dist/**', 'build/**', 'node_modules/**', '**/*.d.ts'],
  },
  js.configs.recommended,
  {
    files: ['**/*.{js,jsx,ts,tsx}'],
    plugins: {
      react,
      'react-hooks': reactHooks,
    },
    languageOptions: {
      parser: require('@babel/eslint-parser'),
      parserOptions: {
        requireConfigFile: false,
        babelOptions: {
          presets: ['@babel/preset-react', '@babel/preset-typescript'],
        },
        ecmaFeatures: { jsx: true },
      },
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    settings: {
      react: {
        version: 'detect',
      },
    },
    rules: {
      ...react.configs.recommended.rules,
      ...reactHooks.configs.recommended.rules,
      'react/react-in-jsx-scope': 'off',
      'react/prop-types': 'off',
      'no-unused-vars': 'warn',
      'no-undef': 'off',
    },
  },
  prettier,
];
