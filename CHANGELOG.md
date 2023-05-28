<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## Unreleased

### Features

- Add `acvnt` metadata
- Add `asoult` metadata
- Update inflation distribution to 85% staking rewards, 5% usage incentives, 10% community pool
- Update signed blocks window to 5000
- Update evm tx base fee to 0.1*1e9 CVN
- Add `SOUL` token contract

### Bug Fixes

- Fix remove tx flags from root cmd
- Fix `cvnd query epochs epoch-infos` use `clientCtx.PrintProto()` print result

## [v1.0.2] - 2023-06-09

### Improvements

- (deps) Bump SDK to v0.46.13

### Bug Fixes

- (vesting) Apply ClawbackVestingAccount Barberry patch

## [v1.0.1] - 2023-05-28

### Improvements

- (deps) Bump IBC-go version to [`v6.1.1`](https://github.com/cosmos/ibc-go/releases/tag/v6.1.1)

### Bug Fixes

- (deps) Bump cosmos-sdk version to `v0.46.10-ledger.3`. 
  Fix memory leak in `cosmos/iavl` package.
- (rpc) [#1431](https://github.com/evmos/evmos/pull/1431) Fix websocket connection id parsing
- Fix math.MaxUint32 overflows int when build cvnd-arm64