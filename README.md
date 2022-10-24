# `ghx`: GitHub Extras (WIP)

Module `ghx` provides packages that expose Go APIs for features of
GitHub that currently do not have a public API:
* Explore
* Topics
* Trending
    * Repositories
    * Developers
* Languages
    * List

Additionally, `ghx` has packages that expose higher-level APIs for convenience:
* [`ghx/stars`](stars): exposes (bulk) APIs related to starring repositories:
    * Star repositories from:
        * Users
        * Organizations
        * Arbitrary URLs that link to GitHub repositories
    * Unstar repositories: 
        * From users
        * From organizations
        * That are archived
        * That have not been udpated in a specified time period