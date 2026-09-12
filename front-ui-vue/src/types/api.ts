export interface ApiResponse<T> {
  code: number
  message: string
  data: T
  request_id: string
}

export interface UserInfo {
  id: number
  identity: string
  name: string
  avatar: string
  email: string
  phone: string
  signature: string
  birth_date?: string
  created_at: string
}

export interface LoginResult {
  token: string
  user: UserInfo
}

export interface UserDetail extends UserInfo {
  is_friend: boolean
  is_self: boolean
}

export interface UserSearchItem {
  id: number
  identity: string
  name: string
  avatar: string
  is_friend: boolean
}

export interface FriendItem {
  id: number
  name: string
  avatar: string
  online: boolean
}

export interface GroupItem {
  id: number
  name: string
  avatar: string
  owner_id: number
  member_count: number
}

export interface OutgoingMessage {
  id: string
  conversation_id: string
  sender_id: number
  sender_name: string
  sender_avatar: string
  content: string
  content_type: 'text' | 'image'
  created_at: string
}

export interface RecentConversation {
  conversation_id: string
  type: 'private' | 'group'
  peer_id: number
  peer_name: string
  peer_avatar: string
  online: boolean
  last_message_content: string
  last_sender_name: string
  last_time: string
}

export interface CaptchaResult {
  captcha_id: string
  image: string
}

export interface UploadResult {
  url: string
}
