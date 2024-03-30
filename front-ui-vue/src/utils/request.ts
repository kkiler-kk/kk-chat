import axios from 'axios';
import qs from 'qs';
import { message, Modal } from 'ant-design-vue';
import { Session } from '@/utils/storage';
const service = axios.create({
    // 环境变量，需要在.env文件中配置
    baseURL: import.meta.env.VITE_API_URL,
    // 超时时间暂定5s
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
    paramsSerializer: {
        serialize(params) {
            return qs.stringify(params, { allowDots: true, arrayFormat: 'brackets' });
        },
    },
});


// 添加请求拦截器
service.interceptors.request.use(
    (config) => {
        // 在发送请求之前做些什么 token
        if (Session.get('token')) {
            config.headers!['Authorization'] = `Bearer ${Session.get('token')}`;
        }
        if (config.method === "get") {
            // get 传递的数组进行处理
            config.paramsSerializer = function (params: any) {
                return qs.stringify(params, { arrayFormat: "repeat" });
            };
        }
        return config;
    },
    (error) => {
        // 对请求错误做些什么
        return Promise.reject(error);
    }
);
// 添加响应拦截器
service.interceptors.response.use(
    (response) => {
        // 对响应数据做点什么
        const { status } = response;
        const res = response.data;
        const code = response.data.code
        if (response.headers["reset-authorization"]) {
            Session.set('token', response.headers["reset-authorization"]);
        }
        if (status == 200) {
            const { code, data, msg } = response.data;
            if (code === 200) {
                return Promise.resolve(data);
            } else {
                message.error(msg);
                return Promise.reject(data);
            }
        } else if (status == 401) {
            if (code === 401) {
                Modal.warning({
                    title: () => '提示',
                    content: () => '登录状态已过期，请重新登录'
                })
                Session.clear(); // 清除浏览器全部临时缓存
                window.location.href = '/'; // 去登录页
            } else if (code !== 0) {
                message.error(res.message)
                return Promise.reject(new Error(res.message))
            } else {
                return res
            }
            return Promise.reject(response)
        }
        return Promise.reject(response)
    },
    (error) => {
        // 对响应错误做点什么
        if (error.message.indexOf('timeout') != -1) {
            message.error('网络超时');
        } else if (error.message == 'Network Error') {
            message.error('网络连接错误');
        } else {
            if (error.response.data) message.error(error.response.statusText);
            else message.error('接口路径找不到');
        }
        return Promise.reject(error);
    }
);
// 导出 axios 实例
export default service;
