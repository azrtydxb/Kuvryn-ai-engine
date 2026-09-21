# Workflow rules

Repo-level rules the procoder skills read and follow. Edit freely — what is written here wins over the skills' built-in defaults.

## Worktrees

Feature work happens in a git worktree when parallel writers would touch overlapping files or branches. Use the harness's native worktree support when available, `git worktree add` otherwise.

## Merge watching

While a PR's checks and reviews run, delegate waiting to a background watcher when available. Watchers only gather information; the main agent owns fixes, replies, and merges.

## After a successful merge

Clean up merged branches and temporary worktrees, run `git fetch --prune`, and return to an updated default branch.
