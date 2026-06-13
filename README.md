# rollops-plugin-growthbook

A [Rollops](https://github.com/klarlabs-studio/rollops) feature-flag provider
plugin backed by [GrowthBook](https://www.growthbook.io/). It drives a feature's
enabled state and a single rollout rule's coverage to track a Rollops canary in
lockstep — as a rollout steps 10% → 50% → 100%, the flag's coverage follows
(0.1 → 0.5 → 1.0).

## How it works

Rollops calls the plugin per progressive step (and/or on promote) with the flag
key, the target environment, and the current traffic percentage. The plugin
POSTs a per-environment update so the environment is enabled/disabled and
carries one `rollout` rule whose `coverage` equals the percentage. GrowthBook
merges the per-environment patch, leaving other environments untouched.

## Configuration

Credentials come from the plugin's own environment, never from the Rollops
target spec:

| Env var              | Required | Default                     | Description               |
|----------------------|----------|-----------------------------|---------------------------|
| `GROWTHBOOK_API_URL` | no       | `https://api.growthbook.io` | Base URL (self-hosted ok) |
| `GROWTHBOOK_TOKEN`   | yes      | —                           | API secret (`Bearer`)     |

## Install

```sh
rollops plugin install growthbook
```

Or build and pin manually with `make build` / `make checksum`, then wire into a
rollout spec:

```yaml
featureFlags:
  plugin: ~/.rollops/plugins/growthbook
  sha256: <pin>
  flag: checkout
  environment: production
  when: both
```

## License

MIT
