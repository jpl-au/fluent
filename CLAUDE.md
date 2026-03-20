# Claude Code Instructions

- Do not add `Co-Authored-By` lines to commit messages
- Do not commit CLAUDE.md
- Do not commit `*-audit.md` files
- Do not push to remote unless explicitly asked
- Use British spelling in code comments and commit messages
- Use `go build ./...` to check compilation - never `go build .` which writes a binary
- When editing Go files, do not worry about precise whitespace or indentation - `gofmt` will normalise formatting automatically

## Code Principles

All code in this repository should follow these principles:

- Is easy to read from top to bottom
- Does not assume that you already know what it is doing
- Does not assume that you can memorise all of the preceding code
- Does not have unnecessary levels of abstraction
- Does not have names that call attention to something mundane
- Makes the propagation of values and decisions clear to the reader
- Has comments that explain why, not what, the code is doing to avoid future deviation
- Has documentation that stands on its own
- Has useful errors and useful test failures
- May often be mutually exclusive with "clever" code
