# go-refs

Have you ever built a Go package which others depend on and need to understand
how to best approach refactoring to limit breaking your external interface?

Go-refs uses Go's AST package to parse a (go) file and list any identifiers
which are in use from a given import path.

## Usage

```sh
go list -json ./... \
  | jq -r '.Dir + "/" + .GoFiles[]' \
  | xargs -n1 go-refs -pkg github.com/hashicorp/terraform/helper/schema \
  | sort | uniq -c | sort -nr
```
