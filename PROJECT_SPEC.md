# Field Expedition Ledger

## 项目目标

Field Expedition Ledger is a small local service for recording geological field expeditions and the observations collected during them. It keeps the three parts of a field notebook connected: an expedition plan, dated site observations, and specimen custody records.

## 用户角色

- Field lead: creates an expedition, starts field work, and closes the expedition.
- Field observer: records site observations with coordinates, elevation, rock type, and confidence.
- Collection custodian: registers specimens and reviews their custody state.

## 核心实体

- Expedition: the named field trip with a region, lead, lifecycle status, and notes.
- Observation: a dated site reading linked to an expedition.
- Specimen: a collected material sample linked to an expedition and custody state.
- Expedition summary: derived counts and elevation statistics for review.

## 业务流程

1. Create and activate an expedition. The field lead submits a name, region, lead, and start date; the service validates the request, stores a planned expedition, and activates it before field notes can be added.
2. Record a site observation. The field observer chooses an active expedition, submits a site code and measurements, and receives a stored observation with a generated identifier.
3. Register and review a specimen. The collection custodian submits a label, material, weight, and collection date, then retrieves the specimen to confirm custody state.
4. Close an expedition and review its summary. The field lead closes an active expedition, after which the service reports observation count, specimen count, and elevation range.

## 状态与规则

- Expeditions move from planned to active to closed; a closed expedition cannot accept new observations or specimens.
- Names, regions, leads, site codes, labels, and materials are required.
- Latitude must be between -90 and 90, longitude between -180 and 180, and confidence must be between 0 and 1.
- Specimen weight must be positive and an observation elevation cannot be below sea level.
- Dates are RFC3339 strings and are normalized to UTC.

## 接口与验证

- `GET /health` reports service readiness.
- `POST /api/expeditions` creates an expedition.
- `GET /api/expeditions` lists expeditions.
- `POST /api/expeditions/{id}/activate` starts an expedition.
- `POST /api/expeditions/{id}/observations` records an observation.
- `GET /api/expeditions/{id}/observations` lists observations.
- `POST /api/expeditions/{id}/specimens` registers a specimen.
- `GET /api/specimens/{id}` retrieves a specimen.
- `POST /api/expeditions/{id}/close` closes an expedition.
- `GET /api/expeditions/{id}/summary` returns derived counts and elevations.

The service is verified with `go build ./...`, `go test ./...`, and runtime HTTP checks for creation, activation, observation recording, specimen registration, and closing.
