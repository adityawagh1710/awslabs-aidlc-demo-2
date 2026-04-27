# auth-service API

Base path: `/auth`

## POST /auth/register
Register a new user.

**Request**
```json
{ "email": "user@example.com", "password": "password123" }
```
**Response** `201`
```json
{ "access_token": "...", "refresh_token": "..." }
```
**Errors**: `409` email taken, `422` validation error

---

## POST /auth/login
Authenticate and receive tokens.

**Request**
```json
{ "email": "user@example.com", "password": "password123", "mfa_code": "123456" }
```
**Response** `200`
```json
{ "access_token": "...", "refresh_token": "..." }
```
**Errors**: `401` invalid credentials/MFA, `422` mfa_required, `429` account locked

---

## POST /auth/refresh
Rotate refresh token and get new access token.

**Request**
```json
{ "refresh_token": "..." }
```
**Response** `200`
```json
{ "access_token": "...", "refresh_token": "..." }
```
**Errors**: `401` invalid/expired token

---

## POST /auth/logout
`Authorization: Bearer <access_token>` required.

**Request**
```json
{ "refresh_token": "..." }
```
**Response** `204`

---

## POST /auth/mfa/enroll
`Authorization: Bearer <access_token>` required.

**Response** `200`
```json
{ "secret": "...", "qr_url": "otpauth://..." }
```

---

## POST /auth/mfa/verify
`Authorization: Bearer <access_token>` required.

**Request**
```json
{ "code": "123456" }
```
**Response** `200`
**Errors**: `401` invalid code

---

## GET /health
**Response** `200`
```json
{ "status": "ok", "service": "auth-service" }
```
