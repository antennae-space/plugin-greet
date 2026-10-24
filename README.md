# Greet, an example plugin of the Web IDE

The example external plugin of [Web IDE](https://github.com/antennae-space/web-ide): a tool that greets you and counts the files of the project, a tab view, a command, and an assistant tool, with a backend. It is built on the [plugin SDK](https://github.com/antennae-space/plugin-sdk): its views in Solid with the UI kit, its backend on `antennae.space/plugin-sdk/backend`. Fork it to start a plugin of your own; the how-to [Write a plugin](https://github.com/antennae-space/web-ide/blob/main/docs/howto/write-a-plugin.md) explains every part.

## Install it

In the IDE, *Settings → Plugins → Install a plugin*: `https://github.com/antennae-space/plugin-greet`, version empty (the latest). Accept its two permissions, then tick its box on the project: its icon comes in the right rail.

- *Read the files of the project* (`files:read`): its backend may read the project.
- *Call fs.list of the IDE* (`core:fs.list`): its view and its backend may list the files through the IDE.

Untick a permission in *Settings → Plugins* and the count is refused at once.

## What is in it

| Path | What |
|---|---|
| `plugin.json` | Its manifest: id, version, the IDE versions it runs on (`ide`), the address of its archives (`assets`), its backend, its permissions, its tool, tab view, command and assistant tool (with its schema) |
| `web/src/index.tsx` | Its page module: two views in Solid with the SDK's kit and a command, `ide.register({views, commands})`; the IDE runs it in a sandboxed frame |
| `vite.config.js` | The SDK's preset: builds `web/index.js`, one ES module |
| `i18n/fr.json`, `i18n/es.json` | Its texts in French and Spanish (`ide.t`), English being the key |
| `backend/main.go` | Its backend on `backend.Serve` of the SDK, run confined; it counts the entries of the project through the IDE (`backend.Core`, `fs.list`) |
| `scripts/pack.sh` | Builds the page module, then `dist/greet-linux-amd64.tar.gz` and `dist/greet-linux-arm64.tar.gz` |
| `.github/workflows/release.yml` | A tag `v*` publishes the release with those archives |

## Release a version

1. Change `version` in `plugin.json` (`1.0.1`).
2. Commit, tag `v1.0.1`, push the tag: the workflow builds and publishes the archives.

Never move a tag once released: the IDE keeps the hash of each archive it installed and refuses another one under the same version.

## Try a version before releasing it

`npm ci`, then `./scripts/pack.sh` builds the archives locally (`npm run build` alone the page module, `go vet ./...` checks the backend). The IDE downloads archives over https only: push a pre-release tag (`v1.0.1-rc1`, with that version in `plugin.json`) and install it with the constraint `=1.0.1-rc1`; a plain constraint never picks a pre-release.

## License

MIT.
