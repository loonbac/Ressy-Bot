# admin

Serves the Ressy control panel over HTTP as an installable module.

Public API: `admin.New(Config) core.Module`. `Config` carries the Discord OAuth client values, the session secret, the plugin directory, the panel address and the admin user id allowlist.

Allowed imports: the Go standard library and `src/modules/core`. The panel must never import another feature module.

Invariants: `Register` starts an HTTP listener and stores the host command catalog and module settings; `Close` shuts the listener down. OAuth state and session cookies are HMAC-signed with the session secret. A session is accepted only for an id in `AdminUserIDs`.

The plugin loader accepts only a literal `.so` basename already present in `PluginDirectory`; no path traversal and no source compilation from the panel.

To add a panel capability, add a new route in `panel.go` with its own file beside it and a focused `httptest` test. To change the Discord OAuth exchange, edit `oauth.go`.
