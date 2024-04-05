
import request from '@/utils/request'

export function addFriend(data) {
    return request({
        url: 'api/app/friend/add',
        method: 'post',
        data: data,
    })
}

export function listFriend() {
    return request({
        url: 'api/app/friend/list',
        method: 'post',
    })
}