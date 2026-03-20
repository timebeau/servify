## Summary

- Describe the user-facing or engineering change.

## Validation

- [ ] Ran the relevant tests locally
- [ ] Ran repo hygiene checks when touching generated/runtime/build boundaries
- [ ] Updated generated assets or lockfiles if required
- [ ] Updated docs/backlog when changing architecture or workflow

## Repo Hygiene Checklist

- [ ] No runtime artifacts are committed
- [ ] No local binaries such as `server` or `server.exe` are committed
- [ ] Files under `uploads/`, `apps/server/uploads/`, `internal/handlers/uploads/`, or `.runtime/` are not committed
- [ ] Any committed generated files are intentional and documented
