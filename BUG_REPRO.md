# Reproducing the incremental observation query

This service exposes observations for an expedition at:

```text
GET /api/expeditions/{id}/observations?since=<RFC3339 timestamp>
```

To reproduce the issue, create an active expedition with two observations:

1. Record one observation at an RFC3339 timestamp chosen as the sync cursor.
2. Record a second observation at a later timestamp.
3. Request the observations endpoint with `since` set to the first timestamp.

The response must be suitable for incremental synchronization: it should contain
only records later than the supplied cursor, preserving the endpoint's normal
JSON page shape and ordering. A malformed `since` value must continue to use the
existing invalid-input response behavior.

