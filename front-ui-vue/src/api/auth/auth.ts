import request from '@/utils/request'

export function sendEmailCode(eamil: string){
    return request({
        url: '/api/app/auth/send/email/code',
        method: 'post',
        data: {'email': eamil}
    })
}
export function captchaGenerate(){
    return request({
        url: '/api/app/auth/captchaGenerate',
        method: 'post',
    })
}