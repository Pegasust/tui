# Changelog
All notable changes to this project will be documented in this file. See [conventional commits](https://www.conventionalcommits.org/) for commit guidelines.

- - -
## [0.1.1](https://github.com/paisano-nix/tui/compare/0.1.0..0.1.1) - 2023-04-18
#### Bug Fixes
- pass args along to the final invocation - ([83af50d](https://github.com/paisano-nix/tui/commit/83af50d6c058999094bfab633e0a50faedafa1d1)) - [@blaggacao](https://github.com/blaggacao)
- cog config - ([507cd13](https://github.com/paisano-nix/tui/commit/507cd138a26807aac2be5b859c9000aad8283203)) - [@blaggacao](https://github.com/blaggacao)
#### Miscellaneous Chores
- add instructions to publish release notesd - ([ad7ba7a](https://github.com/paisano-nix/tui/commit/ad7ba7a1cbc0302103a25d1262c9eb55c0939223)) - [@blaggacao](https://github.com/blaggacao)

- - -

## [0.1.0](https://github.com/tui/paisano-nix/compare/5eef783baf77df737e33e8265834ac8afd0b78df..0.1.0) - 2023-04-17
#### Bug Fixes
- polish the completion ux and add some bling - ([3859076](https://github.com/tui/paisano-nix/commit/38590763cbbdf3175cf62b1c693f83c449313e54)) - [@blaggacao](https://github.com/blaggacao)
- infinit loop on prj-spec init if outside a project repo - ([9b2bf76](https://github.com/tui/paisano-nix/commit/9b2bf7679b671319b96fc24f534620f1d9f27f0f)) - [@blaggacao](https://github.com/blaggacao)
- branding on the `check` sub command - ([1c84e60](https://github.com/tui/paisano-nix/commit/1c84e604adb8907bc20ee5030bb124020ac79ace)) - [@blaggacao](https://github.com/blaggacao)
- nil deref - damint cobra - ([f0272b3](https://github.com/tui/paisano-nix/commit/f0272b3986fbf153322b6e1c8b13016830e3577a)) - [@blaggacao](https://github.com/blaggacao)
- bump std - ([802958d](https://github.com/tui/paisano-nix/commit/802958d123b0a5437441be0cab1dee487b0ed3eb)) - [@blaggacao](https://github.com/blaggacao)
- oversight so that current system is detected again - ([bf8ef13](https://github.com/tui/paisano-nix/commit/bf8ef13f4ad9c84e7bf177c8a5f1c9586c41a4e4)) - [@blaggacao](https://github.com/blaggacao)
#### Continuous Integration
- add gh pages action - ([9756b9a](https://github.com/tui/paisano-nix/commit/9756b9aacc3ab369016c5b56677bf0e8902e8e01)) - [@blaggacao](https://github.com/blaggacao)
#### Documentation
- add tagline description - ([9f03a91](https://github.com/tui/paisano-nix/commit/9f03a911b9293acd93c3fbb1cf1cdaa92ec89c13)) - [@blaggacao](https://github.com/blaggacao)
- add flake-view for docs with mdbook-paisano-preprocessor - ([f45d054](https://github.com/tui/paisano-nix/commit/f45d054b1329e70e475eb185367d18fa08a6a176)) - [@blaggacao](https://github.com/blaggacao)
- fix intro page link - ([92488a2](https://github.com/tui/paisano-nix/commit/92488a29c7b9feac773feba8672d406d5268e3ae)) - [@blaggacao](https://github.com/blaggacao)
- improve wording - ([0fe8858](https://github.com/tui/paisano-nix/commit/0fe88586963807b918cab3e4a6a651604b0a82c2)) - [@blaggacao](https://github.com/blaggacao)
- add rebranding example - ([f32aaec](https://github.com/tui/paisano-nix/commit/f32aaec2774be698590c45438c2b8d0d5cbfa87e)) - [@blaggacao](https://github.com/blaggacao)
- add documentation - ([6c52cf0](https://github.com/tui/paisano-nix/commit/6c52cf0de2e0acd88aef3515f909936abfebb4b6)) - [@blaggacao](https://github.com/blaggacao)
#### Features
- improve description on CLI completion - ([830d91f](https://github.com/tui/paisano-nix/commit/830d91ff32d3e12a4f89dec2f74179416af513c8)) - [@blaggacao](https://github.com/blaggacao)
- comply with PRJ Spec (akin XDG_*) - ([de2574d](https://github.com/tui/paisano-nix/commit/de2574dc7390a9f38ace10b3cb3b35737595f365)) - [@blaggacao](https://github.com/blaggacao)
- add license - ([cb9ac8b](https://github.com/tui/paisano-nix/commit/cb9ac8bc142c6bfac2bebb6566a03175aeb97a05)) - [@blaggacao](https://github.com/blaggacao)
- actions on current system when (remote) build for other system - ([cd31e1c](https://github.com/tui/paisano-nix/commit/cd31e1c13aa01fa811d21b522215037c57e03cd3)) - [@blaggacao](https://github.com/blaggacao)
#### Miscellaneous Chores
- instrument release - ([40fab50](https://github.com/tui/paisano-nix/commit/40fab501a95f1a7f966f0b392557a01c1bcd2b60)) - [@blaggacao](https://github.com/blaggacao)
- add hint to commit readme files - ([f080910](https://github.com/tui/paisano-nix/commit/f0809101b957e831ff5ae3be432397a0da9149b7)) - [@blaggacao](https://github.com/blaggacao)
#### Refactoring
- improve and clean up the code; optional `nom` support - ([2896332](https://github.com/tui/paisano-nix/commit/2896332e412153d7110bac3ebf330e9c5e34404b)) - [@blaggacao](https://github.com/blaggacao)
- use new and shiny paisano direnv support - ([7db9c76](https://github.com/tui/paisano-nix/commit/7db9c76c3e440a926faf3efa585faf1d080585de)) - [@blaggacao](https://github.com/blaggacao)
- make branding configurable at build time - ([694baa7](https://github.com/tui/paisano-nix/commit/694baa76fd58492b721f9091f2ed6736bfa6d85e)) - [@blaggacao](https://github.com/blaggacao)

- - -

Changelog generated by [cocogitto](https://github.com/cocogitto/cocogitto).