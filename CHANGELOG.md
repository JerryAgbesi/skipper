# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.0]

### Added

- `skipper add <alias> <user@host[:port]>` subcommand for writing a new host
  entry to the SSH config. Inherits the `-c, --config` persistent flag.
- Idempotent re-runs of `skipper add` with the same alias and target now
  report `host "<alias>" already exists ..., no change` instead of claiming
  the entry was added.

### Removed

- **BREAKING**: the `-a, --add` flag on the root command has been removed.
  Migrate to the `add` subcommand:

  ```bash
  # before
  skipper --add devone user@10.0.0.8:9000

  # after
  skipper add devone user@10.0.0.8:9000
  ```

## [0.1.5]

- Previous releases are tracked in [GitHub Releases](https://github.com/JerryAgbesi/skipper/releases).

[0.1.5]: https://github.com/JerryAgbesi/skipper/releases/tag/v0.1.5
