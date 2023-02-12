/*
|-------------------------------------------------------------------------------
| Development config               https://maizzle.com/docs/environments/#local
|-------------------------------------------------------------------------------
|
| The exported object contains the default Maizzle settings for development.
| This is used when you run `maizzle build` or `maizzle serve` and it has
| the fastest build time, since most transformations are disabled.
|
*/

module.exports = {
  build: {
    posthtml: {
      expressions: {
        delimiters: ['[[', ']]'],
        unescapeDelimiters: ['[[[', ']]]'],
      }
    },
    templates: {
      source: 'src/templates',
      destination: {
        path: 'build_local',
      },
    },
  },
  albums: Array(5).fill().map((x,i) => i),
  tags: [
    "rock",
    "alternative",
    "indie.rock",
    "drone",
    "new.york"
  ],
  year: () => new Date().getFullYear(),
}
