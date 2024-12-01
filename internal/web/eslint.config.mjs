import pluginJs from "@eslint/js";
import simpleImportSort from "eslint-plugin-simple-import-sort";
import globals from "globals";

export default [
  {
    ignores: ["eslint.config.mjs"],
    languageOptions: { globals: globals.browser },
    plugins: {
      "simple-import-sort": simpleImportSort,
    },
    rules: {
      strict: ["error", "global"],
      "no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          destructuredArrayIgnorePattern: "^_",
          ignoreRestSiblings: true,
        },
      ],
      "simple-import-sort/imports": "error",
      "simple-import-sort/exports": "error",
      "no-restricted-imports": [
        "error",
        {
          patterns: [
            {
              regex: "^(?!.*\\.(js|mjs|json)$).*",
              message: "Use file extension in imports",
            },
          ],
        },
      ],
    },
  },
  {
    files: ["**/*test.mjs"],
    rules: {
      "no-restricted-imports": "off",
    },
  },
  pluginJs.configs.recommended,
];
