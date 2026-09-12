/** 与后端 apperror 业务错误码一致（spec §3.2） */
export const ErrCode = {
  OK: 0,
  InvalidParam: 10001,
  Internal: 10002,
  CaptchaInvalid: 20001,
  BadCredentials: 20002,
  TokenInvalid: 20003,
  UserBanned: 20004,
  EmailTaken: 30001,
  IdentityTaken: 30002,
  UserNotFound: 30003,
  EmailCodeInvalid: 30004,
  AlreadyFriend: 40001,
  NotFriendLimit: 40002,
  GroupNotFound: 50001,
  NotGroupMember: 50002,
  ConversationInvalid: 60001,
} as const

export type ErrCodeValue = (typeof ErrCode)[keyof typeof ErrCode]
