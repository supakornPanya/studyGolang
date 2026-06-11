# Authentication & Authorization Architecture

This document describes the design, trust model, and implementation details of the authentication and authorization system in this project.

---

## 1. Overview & Trust Model

The system utilizes a **Delegated Identity & Local Session** pattern designed to operate securely behind an upstream gateway/identity proxy (such as OAuth / Active Directory / Parent Gateway).

```
   +------------------+         +------------------+         +--------------------------+
   |  Client / UI     |         | Upstream Gateway |         |   KBank ECMS Backend     |
   +------------------+         +------------------+         +--------------------------+
            |                             |                               |
            |-- 1. Login (OAuth/AD) ------>|                               |
            |<-- 2. Parent Token (JWT) ----|                               |
            |                             |                               |
            |-- 3. POST /token (Parent JWT in Header) -------------------->|
            |                                                             |-- Verify Parent Token
            |                                                             |-- Auto-provision User
            |                                                             |-- Map Roles to Profile
            |                                                             |-- Retrieve Scopes
            |<-- 4. Local Session & Refresh Tokens (JWT) -----------------|
            |                             |                               |
            |-- 5. Protected API Request (Bearer Session JWT) ----------->|
            |                             |                               |-- JWTMiddleware
            |                             |                               |-- UserContextMiddleware
            |                             |                               |-- RequireScope Guard
            |                             |                               |-- CheckOwnership Guard
            |<-- 6. API Response -----------------------------------------|
```

### Trust Split
1. **Upstream Gateway**: Responsible for the cryptographic signature verification, expiration check, and audience validation of the incoming **Parent Token** (the identity provider's token).
2. **ECMS Backend**:
   - Assumes the gateway has successfully verified the parent token's authenticity.
   - Extracts identity claims (username, roles) from the parent token without signature verification to bootstrap local users.
   - Issues and manages its own **Local Session Access & Refresh Tokens**.
   - Enforces granular API authorization scopes and database-level object ownership.

---

## 2. Key Components

| Component / Service | File Path | Primary Responsibility |
| :--- | :--- | :--- |
| **Auth Handler** | [auth_handler.go](file:///c:/KBank-ECMS-Backend/cmd/svc-contstrat-backoffice/handler/auth_handler.go) | Exposes endpoints for session initialization (`POST /token`) and token rotation (`POST /refresh-token`). |
| **Parent JWT Verifier** | [parent_jwt.go](file:///c:/KBank-ECMS-Backend/pkg/auth/parent_jwt.go) | Structurally parses parent identity tokens using unverified decoding to extract claims. |
| **JWT Service** | [jwt.go](file:///c:/KBank-ECMS-Backend/pkg/auth/jwt.go) | Generates, validates, and rotates local-session access and refresh tokens signed via HMAC-SHA256. |
| **JWT Middleware** | [auth_middleware.go](file:///c:/KBank-ECMS-Backend/internal/delivery/http/middleware/auth_middleware.go) | Intercepts requests on protected routes, validates local access tokens, and injects user claims into the context. |
| **User Context Middleware** | [user_context_middleware.go](file:///c:/KBank-ECMS-Backend/internal/delivery/http/middleware/user_context_middleware.go) | Loads the full `User` record (with Profile preloaded) from the database and adds it to the request context. |
| **Scope Guard** | [auth_middleware.go](file:///c:/KBank-ECMS-Backend/internal/delivery/http/middleware/auth_middleware.go) | Restricts API access based on required scopes using OR-semantics. |
| **Ownership Guard** | [ownership.go](file:///c:/KBank-ECMS-Backend/internal/service/ownership.go) | Evaluates resource ownership dynamically to ensure users can only modify/delete their own records unless they hold elevated permissions. |
| **Auditing & DB Hooks** | [base_model.go](file:///c:/KBank-ECMS-Backend/internal/domain/entity/base_model.go) | Automatically sets audit fields (`CreatedBy`, `UpdatedBy`) on GORM queries from the context-bound user ID. |

---

## 3. Database & Mapping Models

Local identity, profile, and authorization models are stored in PostgreSQL.

* **User**: Represents a local user record ([user.go](file:///c:/KBank-ECMS-Backend/internal/domain/entity/user.go)). Linked to a single **Profile**.
* **Profile**: Groups user accounts ([profile.go](file:///c:/KBank-ECMS-Backend/internal/domain/entity/profile.go)). Examples: `CS_Superadmin`, `CS_ITAdmin`, `CS_Viewer`.
* **Permission**: Defines a granular capability consisting of a `SOURCE` (e.g. `decision_rule`, `user`) and an `ACTION` (e.g. `CREATE`, `VIEW_ALL`, `EDIT`, `EDIT_ALL`).
* **Profile Permission**: Junction table (`profile_permissions`) mapping profiles to permitted actions.
* **Login Token History**: Audits login sessions by storing a SHA-256 hash of the parent token and its expiration timestamp.

---

## 4. Key Security Flows

### Flow 1: Session Initialization (`POST /token`)

This public endpoint generates local credentials from a parent identity token.

1. **Header Extraction**: The client passes the parent token via the `X-Access-Token` header.
2. **Parent Token Parsing**: The [ParentJWTVerifier](file:///c:/KBank-ECMS-Backend/pkg/auth/parent_jwt.go) extracts claims.
   - Extracted claims: `username` (email/ID) and `role` (assigned parent roles).
   - Validates that `username` is not empty.
3. **Auto-Provisioning**:
   - Checks the database for a user with the matching username.
   - If missing, it creates a new [User](file:///c:/KBank-ECMS-Backend/internal/domain/entity/user.go) record.
4. **Profile Mapping**:
   - Queries `ProfileRepository` matching the parent token's roles to a local profile code.
   - If no profile matches, it defaults the user to `CS_Viewer` (a safe default profile with minimal permissions).
   - Updates the user's `ProfileID` if it is blank.
5. **Scopes Retrieval**:
   - Calls [GetUserScopes](file:///c:/KBank-ECMS-Backend/internal/repository/permission_repository.go#L26) to query all active permissions linked to the user's profile.
   - Concatenates permissions into `"SOURCE:ACTION"` strings (e.g., `decision_rule:CREATE`, `decision_rule:VIEW_ALL`).
6. **Local JWT Generation**:
   - Signs a local **Access Token** (short duration, e.g., 2 hours) containing user identity, active scopes, and `token_use: "access"`.
   - Signs a local **Refresh Token** (longer duration) containing `token_use: "refresh"`.
7. **Audit Logging**:
   - Hashes the raw parent token using SHA-256 (to protect credentials) and upserts it in `LoginTokenHistory`.
8. **Response**: Returns the tokens, expiration, and granted scopes.

---

### Flow 2: Authenticated API Requests

All protected routes configured in [routes.go](file:///c:/KBank-ECMS-Backend/cmd/svc-contstrat-backoffice/handler/routes.go) pass through two middlewares and optional route guards:

#### Step 2.1: JWT Verification ([JWTMiddleware](file:///c:/KBank-ECMS-Backend/internal/delivery/http/middleware/auth_middleware.go#L78))
1. Checks for a `Bearer` token in the `Authorization` header.
2. Verifies the signature of the token against the local secret key.
3. Ensures `token_use` is exactly `access` to prevent token replay attacks (e.g., using a refresh token to perform actions).
4. Injects `userID` (`uuid.UUID`), `claims` (`*auth.Claims`), and `scopes` (`[]string`) into the request context.

#### Step 2.2: Context Hydration ([UserContextMiddleware](file:///c:/KBank-ECMS-Backend/internal/delivery/http/middleware/user_context_middleware.go#L20))
1. Extracts `userID` from the context.
2. Retrieves the complete user record (including profile details) from the database.
3. Inject the `User` struct into the context under `ctxconsts.UserKey`.

#### Step 2.3: Route-Level Scopes Guard ([RequireScope](file:///c:/KBank-ECMS-Backend/internal/delivery/http/middleware/auth_middleware.go#L177))
* Restricts access to specific API endpoints.
* **OR-Semantics**: Verifies that the client holds at least one of the configured scopes.
* Example:
  ```go
  decisionRules.POST("/import",
      middleware.RequireScope("decision_rule", "CREATE", "EDIT", "EDIT_ALL"),
      wizardHandler.ImportDecisionRules)
  ```
  *(Grants access if the user has `decision_rule:CREATE`, `decision_rule:EDIT`, OR `decision_rule:EDIT_ALL`)*.

#### Step 2.4: Service-Level Ownership Guard ([CheckOwnership](file:///c:/KBank-ECMS-Backend/internal/service/ownership.go#L38))
* Invoked in the business logic layer before mutations (updates/deletions).
* Evaluates if a user has elevated administrative scopes (e.g., `EDIT_ALL` or `DELETE_ALL`).
* **Fast Path (In-Memory)**: If the token contains the elevated scope in context, access is approved immediately.
* **Slow Path (Database Fallback)**: If scopes are absent from the token (older format support), queries the database to check for the elevated permission.
* **Record Owner Check**: If the user lacks elevated permission, the guard compares the user's ID against the record's creator (`BaseModel.CreatedBy`). Access is permitted only if the user is the creator.

---

### Flow 3: Session Refresh (`POST /refresh-token`)

Allows token rotation without re-prompting the parent verification.

1. Extracts the `refresh_token` from the JSON request body.
2. Calls `jwtService.RefreshToken`:
   - Validates the token signature and expiration.
   - Enforces that `token_use` is exactly `refresh`.
3. Mints a new access token and a rotated refresh token carrying the same claims.
4. Returns the new pair; the client discards the old refresh token.

---

## 5. Automatic Audit Stamping

All database entities embed [BaseModel](file:///c:/KBank-ECMS-Backend/internal/domain/entity/base_model.go#L19), which implements GORM hooks to manage audit trails transparently:

* **BeforeCreate**: Stamped with the current time and current user ID for `CreatedAt` and `CreatedBy`. Sets `UpdatedAt` and `UpdatedBy` to match.
* **BeforeSave / BeforeUpdate**: Automatically stamps `UpdatedAt` and `UpdatedBy` with the current time and context-bound user ID.
* **BeforeDelete**: Stamps `UpdatedBy` and `UpdatedAt` before soft-deleting the record.

These hooks resolve the user's ID by invoking `getUserID` on the GORM transaction context, which extracts the value stored under `ctxconsts.UserIDKey`.

---

## 6. Local Development & Testing Bypass Modes

To facilitate local testing and continuous integration without connecting to external active directories:

1. **Scope Bypass (`BYPASS_SCOPE_CHECK`)**:
   - If set to `true`, `1`, or `yes`, the [RequireScope](file:///c:/KBank-ECMS-Backend/internal/delivery/http/middleware/auth_middleware.go#L177) guard skips scope checks, allowing local operations without database permissions setup.
2. **Local Token Bypass (`BYPASS_TOKEN`)**:
   - If `SETENV` is set to `DEVLOCAL` and requests originate from `localhost`, the backend allows requests matching the SHA-256 hash defined in the `BYPASS_TOKEN` env var. This allows testing authenticated API endpoints without calling `/token` first.
3. **Environment Setup**:
   For details on setting up local variables, see [PARENT_JWT_GUIDE.md](file:///c:/KBank-ECMS-Backend/docs/PARENT_JWT_GUIDE.md).
