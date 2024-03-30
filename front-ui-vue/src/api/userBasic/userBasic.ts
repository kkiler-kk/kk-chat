
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