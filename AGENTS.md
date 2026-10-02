# Changing pepic

pepic ships through the kindorg-hq golden path (v3). The conventions — work
item, branch, PR title `type(#N): summary`, green before merge, squash-merge,
"a merge to the default branch ships", reading a PR's delivery state, fix
forward — are in
[kindorg-hq/ci AGENTS.md](https://github.com/kindorg-hq/ci/blob/v3/AGENTS.md);
follow them. How the pipeline works: the ci
[README](https://github.com/kindorg-hq/ci/blob/v3/README.md).

pepic specifics:

- **Default branch is `master`**: branch off it, PRs target it.
- **Tests live in the Dockerfile's `test` stage** (Alpine edge + libvips, CGO:
  `go vet ./... && go test ./...`). Run them with `docker build --target test .`;
  the golden path runs the same stage on every PR and every merge.
- **Keep the image the last stage** of the Dockerfile, so a plain
  `docker build .` yields the runtime image.
- **Go module bumps are `fix(deps)`** (they ship); Actions and base-image bumps
  are `chore(deps)` (`.github/dependabot.yml`).
- Production: https://media.heynik.blog, manifests in kindorg-hq/homelab-k8s
  under `manifests/pepic/`. Re-deliver the latest Release with
  `gh workflow run pipeline.yml --repo kindorg-hq/pepic`.
