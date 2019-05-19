# Design Document

pick is made of two parts:
  - the UI
  - fuzzy searching

All the UI is taken care of [1]

Fuzzy searching needs to be thought through more clearly to deliver the expected
behaviour.

Expected Behaviour:
  - display matched characters and submatch for each displayed candidate
    - e.g when searching "xt" I should get "something/* **x**ota/al**t** */file.ext"
  - results should be ranked by:
    - not sure...
    - smallest submatch

[1] there are some rough edges, e.g retaining the position of a selection where
there are fewer options left after searching some more.
