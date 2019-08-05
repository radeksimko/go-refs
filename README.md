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
```
5313 TypeString
2652 ResourceData
2010 Resource
1947 Schema
 920 Set
 658 TypeList
 597 TypeInt
 561 TypeBool
 524 TypeSet
...
```
