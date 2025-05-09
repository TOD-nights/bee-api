module.exports = {
  root: true,
  parserOptions: {
    parser: '@babel/eslint-parser',
    sourceType: 'module'
  },
  env: {
    browser: true,
    node: true,
    es6: true
  },
  extends: ['plugin:vue/recommended', 'eslint:recommended'],
  rules: {
    "vue/max-attributes-per-line" : 0,
    "vue/no-v-model-argument" : 0,
    "baseIndent":0,
    "vue/html-indent":0,
    "vue/first-attribute-linebreak": 0,
    "vue/html-closing-bracket-newline":0,
    "vue/singleline-html-element-content-newline": 0
  }
}
