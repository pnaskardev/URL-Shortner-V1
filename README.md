# URL-Shortner-V1

## Run All Services Together

```bash
make -j2
```

## Protobuf / Contracts

Proto files live in [`contracts/proto/`](contracts/proto/). The generated Go code
in `contracts/gen/go/` is committed to the repo, so whenever you edit a `.proto`
file you must regenerate it.

Run from the `contracts/` directory:

```bash
cd contracts

# (optional) lint before generating
buf lint

# regenerate Go code into contracts/gen/go/
buf generate
```

The workspace (`go.work`) points all services at the local `contracts` module, so
after `buf generate` the new/changed types are picked up immediately — just rebuild:

```bash
# from repo root
cd contracts && go build ./... && cd ../core && go build ./... && cd ../shortener-service && go build ./...
```

Requires the `buf` CLI and the `protoc-gen-go` plugin on your `PATH`.

---

# System-Design Learning Roadmap

This project doubles as a system-design study vehicle — a URL shortener is the single most common
system-design interview, and this repo is already a microservices + event-driven + contract-first
setup. Each phase below turns a system-design concept into a concrete change in this repo, with the
interview talking points attached. Ordered breadth-first by architectural weight.

**Current reality:** the write path is a half-built event-driven pipeline — the RabbitMQ **producer**
works (publisher confirms + quorum queue), but the **consumer is an empty stub**, nothing is persisted
(`ShortURLDB` is never migrated or written), and there is **no read/redirect path**. So the system does
not yet work end-to-end. The "missing" pieces *are* the curriculum.

**Highest-ROI subset if interviewing soon:** CQRS + cache (P1), outbox + idempotency + KGS (P2),
circuit breaker + rate limit (P4), capacity estimation (P2.4) — ~80% of the interview surface.

> Build order matters: **P0 → P1 first** (you can't demo architecture on a service that panics with no
> read path), and never shard broken data (**P3 after P2**). P4/P5 can interleave. Consider doing each
> phase as its own branch/PR so the git history reads as a system-design portfolio.

## Phase 0 — Make it correct & runnable (DONE)

| # | Fix | The lesson |
|---|-----|-----------|
| 0.1 | Handler read `middlewares.UserIDKey` from **core's** internal package (never set by the microservice middleware → nil panic). Now reads `UserId` from the `InternalShortenRequest` protobuf body. | A service must not import another service's internal packages; the contract is the boundary. |
| 0.2 | Config `mapstricture` typo + `LoadConfig` swallowed unmarshal errors and always returned nil. Now fixed and fail-fast. | Config is a validated contract — fail at startup, not at runtime. |
| 0.3 | *(left intentionally)* Key gen encodes the whole URL → unbounded-length, reversible keys (decode back to the original = info leak). | Redesign in **P2.3 (KGS)**. |

## Phase 1 — Read path, CQRS, caching (highest architectural payoff)

The heart of a URL-shortener interview: **massively read-heavy, write-light, eventually consistent.**

- **1.1 CQRS (command/query split).** Command side exists (`POST /api/shorten` → publish
  `UrlCreatedEvent`, returns 202). Build the **query side**: `GET /{key}` handler + read repository.
  *Talking points:* read:write asymmetry (~100:1+), why writes can be async but reads can't,
  **read-your-own-writes** (the eventual-consistency window between publish and redirect availability).
- **1.2 Consumer / event-driven materialization** — fill `shortener-service/workers/urlShortenerDBWorker.go`:
  consume `url.created`, persist to Postgres. Add `ShortURLDB` to `AutoMigrate` in
  `infrastructure/database/dbHelper.go`; move it from `helpers/views/` into `infrastructure/database/models/`.
  *Talking points:* prefetch/QoS, manual ack after commit (at-least-once), worker concurrency, quorum queue.
- **1.3 Redirect semantics — 301 vs 302.** *Talking points:* 301 permanent → browser/CDN-cached, fewer
  hits but you lose click analytics; 302 → every hit routes to you (needed for P5.2). Opinion: 302 + short TTL.
- **1.4 Cache-aside with Redis** on the `GET /{key}` path (`infrastructure/cache/` dirs already exist).
  *Talking points:* cache-aside vs read/write-through; TTL + **negative caching** (cache "not found" to
  stop cache-penetration on bad keys); whether the consumer warms the cache.

## Phase 2 — Data integrity: outbox, idempotency, KGS (where senior interviews live)

- **2.1 Transactional outbox (dual-write problem).** The write publishes to RabbitMQ with no atomic link
  to a DB write. Fix: write the `short_urls` row **and** an `outbox` row in one Postgres transaction; a
  relay goroutine reads the outbox and publishes via the existing `queue.QueueClient`. *Talking points:*
  dual-write vs outbox vs CDC/Debezium; outbox makes two systems one atomic commit.
- **2.2 Idempotency / effectively-once.** Consumer must survive redelivery: use `UrlCreatedEvent.Id`
  (UUID) as idempotency key with upsert / the `ShortenedURLKey` unique index. *Talking points:*
  "exactly-once delivery is a myth; at-least-once + idempotent consumer = effectively-once."
- **2.3 Key Generation Service (KGS) — the signature question.** Replace the reversible full-URL encoding.
  Build and compare: (a) random base62 + collision retry, (b) counter/ticket server + base62-encode the
  ID, (c) Snowflake IDs. *Talking points:* birthday-paradox collision math, ticket-server SPOF fixed by
  ranges/replicas, guessability/enumeration tradeoff, keyspace sizing (62^7 ≈ 3.5T).
- **2.4 Capacity estimation** (paper). Keyspace + write QPS + read:write ratio + avg URL length →
  storage/yr, QPS, cache working-set, bandwidth. The classic interview opener.

## Phase 3 — Scaling the data tier: sharding, replication, CAP

- **3.1 Read replicas & replication.** Split `dbHelper.go`'s single client into primary (writes) +
  replica (reads); GORM dbresolver does this. *Talking points:* replication lag vs the CQRS consistency
  window, CAP/PACELC (AP-leaning — a redirect tolerates slight staleness for availability + low latency).
- **3.2 Sharding with consistent hashing.** Partition `short_urls` by hash of `ShortenedURLKey`; build a
  shard-router (N instances/schemas). *Talking points:* consistent hashing vs mod-N (rehash storm),
  virtual nodes, hot-shard risk, why the **key** (not user_id) is the shard key, cross-shard query pain.
- **3.3 Schema evolution / contract-first** using `contracts/`: make a backward-compatible change to
  `UrlCreatedEvent` (add a field, keep tag numbers), consume old+new in the worker. *Talking points:*
  forward/backward compat, protobuf field-number discipline, expand/contract DB migrations.

## Phase 4 — Resiliency & traffic management

- **4.1 Circuit breaker + bulkhead** around the shortener call in `core/httpClients/httpClient.go` (has
  retries + backoff + jitter already). *Talking points:* retries alone amplify overload (retry storm —
  backoff+jitter is the mitigation), open/half-open/closed, bulkheads isolate slow dependencies.
- **4.2 Rate limiting** at the core gateway (`// SHOULD BE RATE LIMITED` comments already in both
  `routes.go`). Token bucket per user/IP, Redis-backed. *Talking points:* token bucket vs sliding vs
  fixed window; local vs distributed; enforce at gateway vs service; 429 + Retry-After.
- **4.3 Dead-letter queue + backpressure.** Extend `queue.DeclareQueue()` with a DLX/DLQ; poison messages
  route there after N retries; consumer prefetch = backpressure. *Talking points:* poison-message
  handling, retry-with-backoff vs immediate DLQ, DLQ replay.
- **4.4 API gateway / BFF.** `core` already *is* the gateway (auth, forwarding, retries, `X-Transaction-ID`).
  *Talking points:* gateway vs BFF vs mesh; edge concerns vs business logic; the 2-min `core-service` JWT
  is a stepping stone to mTLS (P5.4).

## Phase 5 — Observability, delivery, frontier

- **5.1 Observability (RED/USE, metrics, tracing).** Have slog JSON + `X-Transaction-ID` + `Server-Timing`
  — extend to Prometheus (RED per route), OpenTelemetry traces propagated core → shortener → RabbitMQ →
  consumer, and `/healthz` + `/readyz` (readiness checks Postgres + RabbitMQ + Redis). *Hard part:*
  trace-context propagation **across the async queue boundary** (inject trace IDs into AMQP headers).
- **5.2 Analytics / click pipeline.** Emit `UrlClickedEvent` from `GET /{key}` to a new queue; a worker
  aggregates. This is *why* 302 over 301. *Talking points:* fire-and-forget off the hot path, batching,
  **HyperLogLog** for approximate unique visitors, stream vs batch.
- **5.3 Delivery & service discovery.** No Dockerfile/compose exists; `core/helpers/utils/constants.go`
  hardcodes `http://localhost:8001`. Add Docker Compose (core, shortener, postgres, rabbitmq, redis),
  replace hardcoded URLs with env-based service names. *Talking points:* DNS discovery vs registries,
  12-factor, health-check-driven routing.
- **5.4 Frontier (design, likely not build):** service mesh / mTLS (the service JWT is the app-layer
  version); cache stampede + hot-key (single-flight / jittered TTL / probabilistic early expiry; local
  tier for celebrity URLs); multi-region / geo (geo-DNS, read-local/write-global). **Saga: a poor fit
  here** — no long multi-step distributed transaction; outbox + idempotency cover consistency (saying "I
  chose *not* to use a saga" shows judgment). **Tests:** zero exist — table-driven unit tests (base62,
  JWT) + a testcontainers integration test.

## Concept → interview-topic coverage

CQRS · event-driven architecture · transactional outbox / dual-write · idempotency & effectively-once ·
KGS / Snowflake / ticket server · capacity estimation · consistent-hashing sharding · replication &
CAP/PACELC · multi-tier caching (stampede, hot-key, penetration) · read-your-own-writes · circuit breaker
/ bulkhead / backpressure · rate limiting · DLQ · API gateway / BFF · service mesh / mTLS · observability
(RED/USE, tracing across async boundaries) · analytics pipeline (HyperLogLog) · schema evolution /
contract-first · service discovery · containerization · multi-region · (saga: deliberately excluded).
