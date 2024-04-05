
import request from '@/utils/request'

export function createUserBasic(data){
    return request({
        url: '/api/app/user/add',
        method: 'post',
        data: data,
    })
}
export function login(data){
    return request({
        url: '/api/app/user/login',
        method: 'post',
        data: data,
    })
}

export function detail(id: string) {
    return request({
        url: '/api/app/user/detail/' + id,
        method: 'post'
    })
}
export function logout() {
    return request({
        url:"/api/app/user/logout",
        method: 'post',
    })
}

export function userSearch(data) {
    return request({
        url:"/api/app/user/search",
        method:"post",
        data: data,
    })
}