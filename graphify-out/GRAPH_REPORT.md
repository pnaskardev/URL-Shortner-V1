# Graph Report - .  (2026-07-04)

## Corpus Check
- Corpus is ~12,983 words - fits in a single context window. You may not need a graph.

## Summary
- 182 nodes · 283 edges · 18 communities (15 shown, 3 thin omitted)
- Extraction: 81% EXTRACTED · 18% INFERRED · 1% AMBIGUOUS · INFERRED: 50 edges (avg confidence: 0.8)
- Token cost: 155,230 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Auth Handlers & Password Security|Auth Handlers & Password Security]]
- [[_COMMUNITY_Target Architecture Diagram|Target Architecture Diagram]]
- [[_COMMUNITY_HTTP Client & Retry Logic|HTTP Client & Retry Logic]]
- [[_COMMUNITY_Shortener Service Bootstrap|Shortener Service Bootstrap]]
- [[_COMMUNITY_Config Singleton & JWT Tokens|Config Singleton & JWT Tokens]]
- [[_COMMUNITY_RabbitMQ Queue Client|RabbitMQ Queue Client]]
- [[_COMMUNITY_URL Shortening Pipeline|URL Shortening Pipeline]]
- [[_COMMUNITY_Build & Debug Config|Build & Debug Config]]
- [[_COMMUNITY_Microservice Auth Middleware|Microservice Auth Middleware]]
- [[_COMMUNITY_Shortener Views & Models|Shortener Views & Models]]
- [[_COMMUNITY_Core Auth Middleware|Core Auth Middleware]]
- [[_COMMUNITY_User Model (Core)|User Model (Core)]]
- [[_COMMUNITY_User Model (Shortener)|User Model (Shortener)]]
- [[_COMMUNITY_Postgres Connection|Postgres Connection]]
- [[_COMMUNITY_Core Module Root|Core Module Root]]
- [[_COMMUNITY_Shortener Module Root|Shortener Module Root]]

## God Nodes (most connected - your core abstractions)
1. `QueueClient` - 16 edges
2. `ApiRouter` - 14 edges
3. `shortener.New (constructor)` - 11 edges
4. `RetryableHTTPClient` - 10 edges
5. `SignInHandler` - 10 edges
6. `ShortenURL` - 10 edges
7. `GetConfig` - 9 edges
8. `BadRequest` - 8 edges
9. `ValidationError` - 8 edges
10. `InternalServerError` - 8 edges

## Surprising Connections (you probably didn't know these)
- `URL-Shortner-V1 Project` --references--> `Core Service`  [INFERRED]
  README.md → .vscode/launch.json
- `URL-Shortner-V1 Project` --references--> `Shortener Service`  [INFERRED]
  README.md → .vscode/launch.json
- `make -j2 Parallel Build` --conceptually_related_to--> `Run All Services Compound Debug Config`  [INFERRED]
  README.md → .vscode/launch.json
- `ApiRouter` --references--> `MicroServiceAuthMiddleware`  [EXTRACTED]
  core/api/routes/routes.go → shortener-service/middlewares/microServiceAuthenticator.go
- `ApiRouter` --calls--> `shortener.New (constructor)`  [EXTRACTED]
  core/api/routes/routes.go → shortener-service/pkg/handlers/shortener/shortenerRepository.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Sign-in authentication flow** — core_pkg_handlers_auth_authrepository_signinhandler, core_helpers_validator_validator_validatestruct, core_helpers_utils_passwordhelper_verifypassword, core_helpers_utils_jwthelper_generateaccesstoken, core_helpers_utils_jwthelper_generaterefreshtoken [INFERRED 0.85]
- **Authenticated request middleware chain** — core_middlewares_authmiddleware_authmiddleware, core_helpers_utils_jwthelper_verifyaccesstoken, core_middlewares_authmiddleware_useridkey, core_pkg_handlers_shortner_shortenerrepo_shortenurl [INFERRED 0.85]
- **Shorten proxy to shortener microservice** — core_pkg_handlers_shortner_shortenerrepo_shortenurl, core_httpclients_httpclient_dowithretry, core_helpers_utils_jwthelper_generatemicroserviceauthtoken, core_helpers_utils_constants_globalconstants [INFERRED 0.85]
- **Async URL shortening write-behind pipeline** — shortener_service_pkg_handlers_shortener_shortenerrepository_shortenurl, shortener_service_helpers_utils_base62helper_encodetobase62, shortener_service_infrastructure_queue_queue_publish, shortener_service_workers_urlshortenerdbworker_worker, shortener_service_helpers_views_shortenerviews_shortenerqueuebody [INFERRED 0.75]
- **RabbitMQ publisher-confirm mechanism** — queue_queueclient, shortener_service_infrastructure_queue_queue_newqueueclient, shortener_service_infrastructure_queue_queue_confirmloop, shortener_service_infrastructure_queue_queue_publish [EXTRACTED 1.00]
- **Run All Services (Core + Shortener)** — debug_run_all_services, svc_core, svc_shortener_service [EXTRACTED 1.00]

## Communities (18 total, 3 thin omitted)

### Community 0 - "Auth Handlers & Password Security"
Cohesion: 0.15
Nodes (21): BadRequest, Ctx, InternalServerError, NotFound, ValidationError, SERVICES type, GenerateRefreshToken, generateSalt (+13 more)

### Community 1 - "Target Architecture Diagram"
Cohesion: 0.09
Nodes (26): Actor / Client, Analytics Service (processes/stores click events: IP, timestamp, user-agent, geo), API Gateway (GraphQL), Auth Middleware (validates jwt/paseto token in auth cookie), Auth Service, Clickhouse/BigQuery (heavy analysis), created worker (email/sms/push on creation), Dashboard Service (+18 more)

### Community 2 - "HTTP Client & Retry Logic"
Cohesion: 0.13
Nodes (17): Repository, Client, ApiRouter, App, DB, Context, DB, auth.New (+9 more)

### Community 3 - "Shortener Service Bootstrap"
Cohesion: 0.16
Nodes (13): Repository, ApiRouter(), App, DB, GetConfig(), Config, LoadConfig(), ConnectToPostgres() (+5 more)

### Community 4 - "Config Singleton & JWT Tokens"
Cohesion: 0.20
Nodes (16): Config, GetConfig, Config, LoadConfig, Config sync.Once singleton rationale, GlobalConstants, GenerateAccessToken, GenerateMicroServiceAuthToken (+8 more)

### Community 5 - "RabbitMQ Queue Client"
Cohesion: 0.17
Nodes (10): Channel, Confirmation, Connection, Error, Mutex, pendingConfirm, QueueClient, confirmLoop (+2 more)

### Community 6 - "URL Shortening Pipeline"
Cohesion: 0.24
Nodes (10): URL_CREATED_QUEUE constant, URL_EXCHANGE constant, EncodeToBase62, ShortenerQueueBody (queue message), ShortenRequest (API body), ShortURLDB (GORM model), DeclareQueue, Publish (confirm-blocking) (+2 more)

### Community 7 - "Build & Debug Config"
Cohesion: 0.38
Nodes (7): make -j2 Parallel Build, Launch Core Debug Config, Launch ShortenerService Debug Config, Run All Services Compound Debug Config, URL-Shortner-V1 Project, Core Service, Shortener Service

### Community 8 - "Microservice Auth Middleware"
Cohesion: 0.29
Nodes (5): MapClaims, Token, VerifyAccessToken(), Handler, MicroServiceAuthMiddleware

### Community 9 - "Shortener Views & Models"
Cohesion: 0.40
Nodes (5): Model, UUID, ShortenRequest, ShortenerQueueBody, ShortURLDB

### Community 10 - "Core Auth Middleware"
Cohesion: 0.40
Nodes (4): AuthMiddleware, Handler, UserIDKey, contextKey

### Community 11 - "User Model (Core)"
Cohesion: 0.50
Nodes (3): Model, User, UUID

### Community 12 - "User Model (Shortener)"
Cohesion: 0.50
Nodes (3): Model, User, UUID

## Ambiguous Edges - Review These
- `ShortenerQueueBody (queue message)` → `URL Shortener DB Worker (stub)`  [AMBIGUOUS]
  shortener-service/workers/urlShortenerDBWorker.go · relation: shares_data_with
- `ShortURLDB (GORM model)` → `URL Shortener DB Worker (stub)`  [AMBIGUOUS]
  shortener-service/workers/urlShortenerDBWorker.go · relation: shares_data_with
- `DeclareQueue` → `URL Shortener DB Worker (stub)`  [AMBIGUOUS]
  shortener-service/workers/urlShortenerDBWorker.go · relation: references

## Knowledge Gaps
- **15 isolated node(s):** `github.com/pnaskardev/URL-Shortner-V1/core`, `contextKey`, `github.com/pnaskardev/URL-Shortner-V1/shortener-service`, `ShortenRequest`, `GlobalConstants` (+10 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `ShortenerQueueBody (queue message)` and `URL Shortener DB Worker (stub)`?**
  _Edge tagged AMBIGUOUS (relation: shares_data_with) - confidence is low._
- **What is the exact relationship between `ShortURLDB (GORM model)` and `URL Shortener DB Worker (stub)`?**
  _Edge tagged AMBIGUOUS (relation: shares_data_with) - confidence is low._
- **What is the exact relationship between `DeclareQueue` and `URL Shortener DB Worker (stub)`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **Why does `ApiRouter` connect `HTTP Client & Retry Logic` to `Auth Handlers & Password Security`, `Shortener Service Bootstrap`, `Config Singleton & JWT Tokens`, `URL Shortening Pipeline`, `Microservice Auth Middleware`, `Core Auth Middleware`?**
  _High betweenness centrality (0.257) - this node is a cross-community bridge._
- **Why does `shortener.New (constructor)` connect `Shortener Service Bootstrap` to `HTTP Client & Retry Logic`, `RabbitMQ Queue Client`, `URL Shortening Pipeline`?**
  _High betweenness centrality (0.150) - this node is a cross-community bridge._
- **Why does `main` connect `Config Singleton & JWT Tokens` to `HTTP Client & Retry Logic`, `RabbitMQ Queue Client`, `Postgres Connection`, `URL Shortening Pipeline`?**
  _High betweenness centrality (0.099) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `shortener.New (constructor)` (e.g. with `ApiRouter()` and `LoadConfig()`) actually correct?**
  _`shortener.New (constructor)` has 5 INFERRED edges - model-reasoned connections that need verification._