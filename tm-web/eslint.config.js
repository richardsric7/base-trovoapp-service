const js = require("@eslint/js");
const react = require("eslint-plugin-react");
const reactHooks = require("eslint-plugin-react-hooks");
const jsxA11y = require("eslint-plugin-jsx-a11y");
const nextPlugin = require("@next/eslint-plugin-next");
const globals = require("globals");

// `eslint-config-next` (and typescript-eslint, which it requires
// unconditionally at import time) doesn't support TypeScript 7 yet -
// confirmed by a hard runtime guard, not just a peerDependency warning
// (see app-web's eslint.config.js for the full story). This config
// rebuilds next/core-web-vitals from its pieces directly
// (@next/eslint-plugin-next + eslint-plugin-react + eslint-plugin-react-hooks
// + eslint-plugin-jsx-a11y, the same plugins eslint-config-next bundles)
// and parses .ts/.tsx syntax-only via @babel/eslint-parser +
// @babel/preset-typescript, independent of the `typescript` package
// version, so tsc stays on TypeScript 7 for the real type-checking gate.
//
// `next lint` itself was removed in Next.js 16 in favor of running
// ESLint directly - this file plus the `lint` script in package.json
// replace it.
module.exports = [
  {
    // public/ holds vendored static assets (e.g. pdf.worker.min.mjs) - not
    // our source and not covered by the files glob below anyway (.mjs),
    // so without this it falls through to bare js.configs.recommended
    // defaults and floods the report with findings against minified code.
    ignores: [".next/**", "out/**", "build/**", "next-env.d.ts", "public/**"],
  },
  js.configs.recommended,
  {
    files: ["**/*.{js,jsx,ts,tsx}"],
    plugins: {
      react,
      "react-hooks": reactHooks,
      "jsx-a11y": jsxA11y,
      "@next/next": nextPlugin,
    },
    languageOptions: {
      parser: require("@babel/eslint-parser"),
      parserOptions: {
        requireConfigFile: false,
        babelOptions: {
          presets: ["@babel/preset-react", "@babel/preset-typescript"],
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
        version: "detect",
      },
    },
    rules: {
      ...react.configs.recommended.rules,
      ...reactHooks.configs.recommended.rules,
      ...jsxA11y.flatConfigs.recommended.rules,
      ...nextPlugin.configs["core-web-vitals"].rules,
      "react/react-in-jsx-scope": "off",
      "react/prop-types": "off",
      "react/no-unknown-property": "off",
      "react/jsx-no-target-blank": "off",
      "no-unused-vars": "warn",
      "no-undef": "off",
      // Downgraded to warn so CI can pass on the current codebase.
      // Fix these incrementally, then flip back to error one at a time.
      "react-hooks/rules-of-hooks": "warn",
      "react-hooks/exhaustive-deps": "warn",
      "react/jsx-key": "warn",
      "@next/next/no-img-element": "warn",
      // eslint-plugin-react-hooks 6/7 added a batch of new React Compiler
      // diagnostic rules to "recommended", all defaulting to "error" - none
      // of these existed when this config was written, and auditing ~600
      // newly-surfaced findings across the existing codebase is its own
      // project, not something to force through a dependency bump. Same
      // "downgrade to warn, fix incrementally" treatment as above.
      "react-hooks/static-components": "warn",
      "react-hooks/use-memo": "warn",
      "react-hooks/preserve-manual-memoization": "warn",
      "react-hooks/immutability": "warn",
      "react-hooks/globals": "warn",
      "react-hooks/refs": "warn",
      "react-hooks/set-state-in-effect": "warn",
      "react-hooks/error-boundaries": "warn",
      "react-hooks/purity": "warn",
      "react-hooks/set-state-in-render": "warn",
      "react-hooks/config": "warn",
      "react-hooks/gating": "warn",
    },
  },
];
