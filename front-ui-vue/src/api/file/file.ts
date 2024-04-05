
import request from '@/utils/request'

export function uploadImages(file) {
    const formData = new FormData();
    formData.append('file', file);
    return request({
        url: '/api/app/file/upload/image',
        method: 'post',
        data: formData,
        headers: {
            'Content-type': 'multipart/form-data',
        },
    })
}