# `ghx`: GitHub Extras (WIP)

Module `ghx` provides packages that expose Go APIs for features of
GitHub that currently do not have a public API:

- Explore
- Topics
- Trending
  - Repositories
  - Developers
- Languages
  - List

Additionally, `ghx` has packages that expose higher-level APIs for convenience:

- [`ghx/stars`](stars): exposes (bulk) APIs related to starring repositories:
  - Star repositories from:
    - Users
    - Organizations
    - Arbitrary URLs that link to GitHub repositories
  - Unstar repositories:
    - From users
    - From organizations
    - That are archived
    - That have not been udpated in a specified time period

# Building

```bash
$ nix build
```

Will produce executable binaries for all discovered `main` packages in
`result/bin/`. Currently those are the 3 `gh` CLI extension compatible:

- `gh-explore` (not implemented yet): CLI wrapper for [`ghx/explore`](explore)
- `gh-languages`: CLI wrapper for [`ghx/languages`](languages)
- `gh-stars` (not implemented yet): CLI wrapper for [`ghx/stars`](stars)
