
import request from '@/utils/request'

export function groupAdd(data) {
    return request({
        url: '/api/app/group/add',
        method: 'post',
        data: data,
    })
}