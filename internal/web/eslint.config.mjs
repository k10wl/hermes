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
  pluginJs.configs.recommended,
];
