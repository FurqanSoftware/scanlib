import { parser } from "./scanspec.grammar"
import {
  LRLanguage,
  LanguageSupport,
  indentNodeProp,
  foldNodeProp,
  foldInside,
} from "@codemirror/language"
import { completeFromList } from "@codemirror/autocomplete"
import { styleTags, tags as t } from "@lezer/highlight"

const parserWithMetadata = parser.configure({
  props: [
    styleTags({
      "var check scan scanln for if else end eol eof": t.keyword,
      TypeName: t.typeName,
      VariableName: t.variableName,
      "Float Integer": t.number,
      String: t.string,
      LineComment: t.lineComment,
      "CompareOp ArithOp MulOp PowerOp Or And": t.operator,
      '"=" AssignOp': t.operator,
      Ellipsis: t.punctuation,
      '"(" ")"': t.paren,
      '"[" "]"': t.squareBracket,
      '","': t.separator,
      '"."': t.punctuation,
    }),
    indentNodeProp.add({
      IfStatement: (cx) => cx.baseIndent + cx.unit,
      ForStatement: (cx) => cx.baseIndent + cx.unit,
    }),
    foldNodeProp.add({
      IfStatement: foldInside,
      ForStatement: foldInside,
    }),
  ],
})

const scanspecLanguage = LRLanguage.define({
  parser: parserWithMetadata,
  languageData: {
    commentTokens: { line: "#" },
  },
})

const keywords = [
  "var", "scan", "scanln", "check",
  "if", "else", "end",
  "for", "eol", "eof",
].map((kw) => ({ label: kw, type: "keyword" }))

const types = [
  "bool", "int", "int64",
  "float32", "float64", "string",
].map((t) => ({ label: t, type: "type" }))

const builtinFunctions = [
  { label: "len", type: "function", detail: "(a)", info: "Returns the length of array a" },
  { label: "toInt64", type: "function", detail: "(s, b=10)", info: "Parses string s in base b and returns int64" },
]

const moduleFunctions = [
  // arrays
  { label: "arrays.sorted", type: "function", detail: "(A)", info: "Returns true if the array is sorted" },
  { label: "arrays.distinct", type: "function", detail: "(A)", info: "Returns true if all elements are unique" },
  { label: "arrays.permutation", type: "function", detail: "(N, A)", info: "Returns true if A is a permutation of 1..N" },
  { label: "arrays.range", type: "function", detail: "(A, lo, hi)", info: "Returns true if all elements are in [lo, hi]" },
  // graphs
  { label: "graphs.simple", type: "function", detail: "(N, U, V)", info: "Returns true if graph has no self-loops or duplicate edges" },
  { label: "graphs.connected", type: "function", detail: "(N, U, V)", info: "Returns true if all N nodes are connected" },
  { label: "graphs.acyclic", type: "function", detail: "(N, U, V)", info: "Returns true if graph has no cycles" },
  { label: "graphs.tree", type: "function", detail: "(N, U, V)", info: "Returns true if edges form a tree on N nodes" },
  // math
  { label: "math.abs", type: "function", detail: "(n)", info: "Returns the absolute value of n" },
  { label: "math.min", type: "function", detail: "(a, b)", info: "Returns the minimum of a and b" },
  { label: "math.max", type: "function", detail: "(a, b)", info: "Returns the maximum of a and b" },
  { label: "math.sqrt", type: "function", detail: "(n)", info: "Returns the square root of n" },
  { label: "math.log2", type: "function", detail: "(n)", info: "Returns the base-2 logarithm of n" },
  { label: "math.ceil", type: "function", detail: "(n)", info: "Returns the smallest integer not less than n" },
  { label: "math.floor", type: "function", detail: "(n)", info: "Returns the largest integer not greater than n" },
  { label: "math.pow", type: "function", detail: "(n, e)", info: "Returns n raised to the power of e" },
  { label: "math.sum", type: "function", detail: "(a...)", info: "Returns the sum of the arguments" },
  // matrices
  { label: "matrices.dimensions", type: "function", detail: "(G, R, C)", info: "Returns true if G has R rows, each of length C" },
  // numbers
  { label: "numbers.prime", type: "function", detail: "(n)", info: "Returns true if n is prime" },
  { label: "numbers.gcd", type: "function", detail: "(a, b)", info: "Returns the greatest common divisor" },
  { label: "numbers.lcm", type: "function", detail: "(a, b)", info: "Returns the least common multiple" },
  { label: "numbers.coprime", type: "function", detail: "(a, b)", info: "Returns true if a and b are coprime" },
  // regexp
  { label: "regexp.match", type: "function", detail: "(s, pattern)", info: "Returns true if s matches the pattern" },
  { label: "regexp.fullmatch", type: "function", detail: "(s, pattern)", info: "Returns true if entire string matches" },
  { label: "regexp.count", type: "function", detail: "(s, pattern)", info: "Returns number of non-overlapping matches" },
  // strings
  { label: "strings.distinct", type: "function", detail: "(s)", info: "Returns true if all characters are unique" },
  { label: "strings.sorted", type: "function", detail: "(s)", info: "Returns true if s is sorted" },
  { label: "strings.palindrome", type: "function", detail: "(s)", info: "Returns true if s is a palindrome" },
  { label: "strings.lowercase", type: "function", detail: "(s)", info: "Returns true if s is all lowercase" },
  { label: "strings.uppercase", type: "function", detail: "(s)", info: "Returns true if s is all uppercase" },
  { label: "strings.alpha", type: "function", detail: "(s)", info: "Returns true if s is all letters" },
  { label: "strings.digit", type: "function", detail: "(s)", info: "Returns true if s is all digits" },
  { label: "strings.alphanumeric", type: "function", detail: "(s)", info: "Returns true if s is all letters and digits" },
  { label: "strings.binary", type: "function", detail: "(s)", info: "Returns true if s is only '0' and '1'" },
  { label: "strings.contains", type: "function", detail: "(s, sub)", info: "Returns true if s contains substring sub" },
]

const staticCompletion = completeFromList([
  ...keywords,
  ...types,
  ...builtinFunctions,
  ...moduleFunctions,
])

function variableCompletion(context) {
  const word = context.matchBefore(/[a-zA-Z_]\w*/)
  if (!word) return null

  const text = context.state.doc.toString()
  const variables = new Set()

  for (const match of text.matchAll(/\bvar\s+([a-zA-Z_]\w*(?:\s*,\s*[a-zA-Z_]\w*)*)/g)) {
    for (const name of match[1].split(",")) {
      variables.add(name.trim())
    }
  }

  for (const match of text.matchAll(/\bfor\s+([a-zA-Z_]\w*)\s*:=/g)) {
    variables.add(match[1])
  }

  return {
    from: word.from,
    options: Array.from(variables).map((v) => ({ label: v, type: "variable" })),
  }
}

export function scanspec() {
  return new LanguageSupport(scanspecLanguage, [
    scanspecLanguage.data.of({
      autocomplete: staticCompletion,
    }),
    scanspecLanguage.data.of({
      autocomplete: variableCompletion,
    }),
  ])
}
