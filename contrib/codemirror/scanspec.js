import { StreamLanguage } from "@codemirror/language"

export const scanspec = StreamLanguage.define({
  token(stream, state) {
    // Comments
    if (stream.match("#")) {
      stream.skipToEnd()
      return "comment"
    }

    // Strings
    if (stream.match('"')) {
      while (!stream.eol()) {
        if (stream.next() === "\\") stream.next()
        else if (stream.current().endsWith('"') && stream.current().length > 1) break
      }
      return "string"
    }

    // Numbers (float before int)
    if (stream.match(/\d+\.\d*/)) return "number"
    if (stream.match(/\d+/)) return "number"

    // Operators
    if (stream.match(/\*\*|==|!=|<=|>=|&&|\|\|/)) return "operator"
    if (stream.match(/[+\-*/<>=!]/)) return "operator"

    // Keywords, types, and identifiers
    if (stream.match(/[a-zA-Z_][a-zA-Z0-9_]*/)) {
      const word = stream.current()
      if (/^(check|else|end|eof|eol|for|if|scan|scanln|var)$/.test(word))
        return "keyword"
      if (/^(bool|int|int64|float32|float64|string)$/.test(word))
        return "typeName"
      return "variableName"
    }

    // Skip whitespace and punctuation
    stream.next()
    return null
  },
})
