# Vox Loop

Private Matrix/Dendrite communication stack on Tailscale for friends & family.

## Architecture

- **Dendrite** — Matrix homeserver (monolith mode)
- **Sliding Sync** — MSC3575 proxy for Element X clients
- **Caddy** — Reverse proxy with automatic TLS
- **Postgres** — Database backend
- **vox-loop CLI** — Configuration generator and container entrypoint

## Quick Start

```bash
# 1. Generate configuration
go run ./cmd/vox-loop init

# 2. Copy and edit environment variables
cp .env.example .env
# Edit .env with real passwords

# 3. Start the stack
docker compose up -d

# 4. Create your first account
docker compose exec dendrite vox-loop admin create-account --username yourname --admin
```

## Roadmap

### Phase 1: The Vox Network

Text communications over Tailscale MagicDNS. Dendrite monolith with Postgres, Sliding Sync proxy for Element X, and Caddy for TLS termination. All services communicate over the Tailscale mesh — no public internet exposure.

- Dendrite homeserver at `imperial-construct.tail64150e.ts.net`
- Element X on mobile and desktop via Tailscale
- Registration locked down; accounts created via CLI
- Sliding Sync for fast room list and message sync

### Phase 2: Tactical Auspex

Voice and video via LiveKit + MatrixRTC. Hardware-accelerated transcoding on the 3080 for media processing. Encrypted voice channels with spatial audio support.

- LiveKit server with MatrixRTC integration
- NVIDIA 3080 transcoding pipeline
- Element Call for encrypted group calls
- Push notifications via UnifiedPush

### Phase 3: The Imperial Failover

High availability and disaster recovery. Postgres replication via Litestream for continuous backup. Edge proxy for optional external federation.

- Postgres WAL streaming to object storage
- Automated failover and recovery
- Edge proxy for federation (optional)
- Monitoring and alerting stack
