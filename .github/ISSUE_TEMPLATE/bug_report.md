---
name: Bug report
about: Report a problem in the backend or webapp
title: "Bug: <short description>"
labels: bug
assignees: ''
---

## Summary
Clear description of the bug.

## Steps to Reproduce
1. …
2. …

Example backend:
```bash
make build && PD_MONGO_URI="mongodb://localhost:27017" ./dist/pd.exe
```
Example webapp:
```bash
cd webapp/_webapp && npm install && npm run build:prd:chrome
```

## Expected vs Actual
- Expected: …
- Actual: …

## Environment
- OS: …
- Go: `go version` → …
- Node: `node -v` → …
- Browser/Extension version: …
- Backend commit: `git rev-parse --short HEAD`

## Logs / Screenshots
- Backend logs (redact secrets)
- Browser console / network errors

## Scope
Affected modules (check): `backend` `webapp/_webapp` `proto` `helm-chart` other: …

## Additional Context
Anything else helpful.

