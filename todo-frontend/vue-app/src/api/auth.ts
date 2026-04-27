import api from './client'

export interface TokenPair { access_token: string; refresh_token: string }
export interface User { id: string; email: string; mfa_enabled: boolean }

export const authApi = {
  register: (email: string, password: string) =>
    api.post<TokenPair & { user: User }>('/auth/register', { email, password }),

  login: (email: string, password: string, mfa_code?: string) =>
    api.post<TokenPair & { user: User; mfa_required?: boolean }>('/auth/login', { email, password, mfa_code }),

  logout: (refresh_token: string) =>
    api.post('/auth/logout', { refresh_token }),

  enrollMfa: () =>
    api.post<{ secret: string; qr_url: string }>('/auth/mfa/enroll'),

  verifyMfa: (code: string) =>
    api.post('/auth/mfa/verify', { code }),
}
