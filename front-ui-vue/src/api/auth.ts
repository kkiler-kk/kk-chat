import { request } from './client'
import type { CaptchaResult, LoginResult } from '@/types/api'

export interface RegisterParams {
  identity: string
  name: string
  password: string
  email: string
  email_code: string
}

export interface LoginParams {
  account: string
  password: string
  captcha_id: string
  captcha_answer: string
}

export const register = (p: RegisterParams) =>
  request<null>({ url: '/api/v1/auth/register', method: 'POST', data: p })

export const login = (p: LoginParams) =>
  request<LoginResult>({ url: '/api/v1/auth/login', method: 'POST', data: p })

export const logout = () => request<null>({ url: '/api/v1/auth/logout', method: 'POST' })

export const getCaptcha = () =>
  request<CaptchaResult>({ url: '/api/v1/auth/captcha', method: 'POST' })

export const sendEmailCode = (email: string) =>
  request<null>({ url: '/api/v1/auth/email-code', method: 'POST', data: { email } })
