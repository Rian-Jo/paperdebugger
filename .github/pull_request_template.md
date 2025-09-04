## Summary
Describe the goal and context. What problem does this PR solve?

## Linked Issues
Closes #<issue>, relates to #<issue>

## Changes
- Brief bullet list of key changes
- Note affected modules: `backend`, `webapp/_webapp`, `proto`, etc.

## Screenshots / GIFs (UI)
If the extension UI changes, include before/after visuals.

## How to Test
- Backend: `make gen fmt lint test` and `make build && PD_MONGO_URI="mongodb://localhost:27017" ./dist/pd.exe`
- Webapp: `cd webapp/_webapp && npm install && npm run build:prd:chrome`

## Breaking Changes
List any API/schema/behavior changes and migration steps.

## Security / Config Notes
- Env vars only (no secrets in code). Update `.env.example` if needed.
- For custom endpoints, ensure HTTPS when used by the browser extension.

## Checklist
- [ ] PR title uses Conventional Commits (e.g., `feat:`, `fix:`, `docs:`)
- [ ] Linked related issues (`Closes #...`)
- [ ] Ran `make gen fmt lint test` successfully
- [ ] Updated docs/README/AGENTS.md as needed
- [ ] Added/updated tests where appropriate
- [ ] No secrets or private data committed
- [ ] UI changes include screenshots/GIFs (if applicable)
