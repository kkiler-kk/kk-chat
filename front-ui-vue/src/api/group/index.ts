
import request from '@/utils/request'

export function groupAdd(data) {
    return request({
        url: '/api/app/group/add',
        method: 'post',
        data: data,
    })
}

export function findByName(name: string) {
    return request({
        url: '/api/app/group/find/' +name,
        method: "post",
    })
}

export function listGroup() {
    return request({
        url: '/api/app/group/list',
        method: 'post',
    })
}