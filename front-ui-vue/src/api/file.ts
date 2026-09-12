import { request } from './client'
import type { UploadResult } from '@/types/api'

export const uploadImage = (file: File) => {
  const form = new FormData()
  form.append('file', file)
  return request<UploadResult>({
    url: '/api/v1/files/images',
    method: 'POST',
    data: form,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
