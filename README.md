# Vox Loop

Private Matrix/Dendrite communication stack on Tailscale for friends & family.

## Architecture

- **Dendrite** — Matrix homeserver (monolith mode)
- **Sliding Sync** — MSC3575 proxy required by Element X
- **Caddy** — Reverse proxy, serves well-known discovery and routes Matrix traffic
- **Postgres** — Database backend
- **vox-loop CLI** — Configuration generator and container entrypoint
- **Tailscale** — Encrypted mesh networking (TLS at the network layer)

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

## Connecting

All users must be connected to Tailscale to reach the homeserver. Traffic is encrypted end-to-end at the network layer by Tailscale, so HTTPS is not required.

### Element X (Recommended)

Element X is the modern Matrix client with native Sliding Sync support for fast sync performance.

- **iOS**: [App Store](https://apps.apple.com/app/element-x-secure-messenger/id6448611190)
- **Android**: [Google Play](https://play.google.com/store/apps/details?id=io.element.android.x)

On the login screen, tap **Change homeserver** and enter:
```
http://imperial-construct:4000
```

### Element (Desktop & Web)

Element is the full-featured Matrix client for desktop and browser.

- **Desktop** (macOS/Windows/Linux): [element.io/download](https://element.io/download)
- **Web**: [app.element.io](https://app.element.io)

Click **Sign in**, then **Edit** the homeserver URL to:
```
http://imperial-construct:4000
```

> **Note:** Element Web at `app.element.io` is served over HTTPS and may block connections to an HTTP homeserver due to mixed-content restrictions. Use the desktop app or Element X instead.

## Roadmap

### Phase 1: The Vox Network

Text communications over Tailscale MagicDNS. Dendrite monolith with Postgres, all traffic encrypted by Tailscale — no public internet exposure.

- Dendrite homeserver at `imperial-construct:4000`
- Sliding Sync proxy for Element X mobile clients
- Element X on mobile and desktop via Tailscale
- Registration locked down; accounts created via CLI

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
