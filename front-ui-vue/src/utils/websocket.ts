import { Local } from "./storage";
import { isJsonString } from '@/utils/is';
let socket: WebSocket;
let isActive: boolean;
export function getSocket(): WebSocket {
    if (socket === undefined) {
        location.reload();
    }
    return socket;
}

export function getActive(): boolean {
    return isActive;
}

export function sendMsg(event: string, data = null, isRetry = true) {
    if (socket === undefined || !isActive) {
        if (!isRetry) {
            console.log('socket连接异常，发送失败！');
            return;
        }
        console.log('socket连接异常，等待重试..');
        setTimeout(function () {
            sendMsg(event, data);
        }, 200);
        return;
    }

    try {
        socket.send(
            JSON.stringify({
                event: event,
                data: data,
            })
        );
    } catch (err) {
        // @ts-ignore
        console.log('ws发送消息失败，等待重试，err：' + err.message);
        if (!isRetry) {
            return;
        }
        setTimeout(function () {
            sendMsg(event, data);
        }, 100);
    }
}
let lockReconnect = false;
let timer: ReturnType<typeof setTimeout>;
export function createSocket() {
    console.log('createSocket...') // 创建websocket
    try {
        if (Local.get('token') == '') {
            throw new Error('用户尚未登录， 请稍后重试。。。')
        }
        socket = new WebSocket(import.meta.env.VITE_WEBSOCKET_URL + "?token=" + Local.get("token"))
        init()
    } catch (e) {
        console.log('createSocket err:' + e);
        reconnect(); // 重连一下websocket
    }

};
export async function onMessage(event) {
    if (socket == undefined || !isActive) {
        console.log("socket连接异常，发送失败!")
        return
    }
    try {
        isActive = true;
        console.log('WebSocket:收到一条消息', event.data);

        let isHeart = false;
        if (!isJsonString(event.data)) {
            console.log('socket message incorrect format:' + JSON.stringify(event));
            return;
        }

        const message = JSON.parse(event.data);
        if (message.event === 'ping') {
            isHeart = true;
        }

        // 强制退出
        if (message.event === 'kick') {
            Session.clear()
            location.reload();
        }

        // 通知
        if (message.event === 'notice') {

        }
        heartCheck.reset().start()

        window.dispatchEvent(new CustomEvent("onmessageWS", {
            detail: {
                event: message
            }
        }))
    } catch (err) {
        return null
    }
    return null
}
export function init() {
    socket.onopen = function (_) {
        console.log('WebSocket:已连接');
        isActive = true;
        //心跳检测重置
        heartCheck.reset().start();
        lockReconnect = true
    };

    socket.onmessage = function (event) {
        isActive = true;
        // console.log('WebSocket:收到一条消息', event.data);

        let isHeart = false;
        if (!isJsonString(event.data)) {
            console.log('socket message incorrect format:' + JSON.stringify(event));
            return;
        }

        const message = JSON.parse(event.data);
        if (message.event === 'ping') {
            isHeart = true;
        }

        // 强制退出
        if (message.event === 'kick') {
            Session.clear()
            location.reload();
        }

        // 通知
        if (message.event === 'notice') {

        }

        window.dispatchEvent(new CustomEvent("onmessageWS", {
            detail: {
                event: message
            }
        }))
        heartCheck.reset().start()
    }

    socket.onerror = function (_) {
        console.log('WebSocket:发生错误');
        reconnect();
        isActive = false;
    };

    socket.onclose = function (_) {
        console.log('WebSocket:已关闭');
        heartCheck.reset()
        reconnect();
        isActive = false;
    };

    window.onbeforeunload = function () {
        socket.close();
        isActive = false;
    };
};
const reconnect = () => {
    if (lockReconnect) return;
    lockReconnect = true;
    timer && clearTimeout(timer);
    timer = setTimeout(() => {
        createSocket();
        lockReconnect = false
    }, 5000);
}

export function addOnMessage(onMessageList: any, func: Function) {
    let exist = false;
    for (let i = 0; i < onMessageList.length; i++) {
        if (onMessageList[i].name == func.name) {
            onMessageList[i] = func;
            exist = true;
        }
    }
    if (!exist) {
        onMessageList.push(func);
    }
}
// const heartCheck = {
//     timeout: 3000,
//     timeoutObj: null,
//     serverTimeoutObj: null,
//     start: function() {
//         var self = this;
//         this.timeoutObj && clearTimeout(this.timeoutObj)
//         this.serverTimeoutObj && clearTimeout(this.serverTimeoutObj)
//         this.timeoutObj = setTimeout(function(){
//             //这里发送一个心跳，后端收到后，返回一个心跳消息，
//             sendMsg("ping",null)
//             // self.serverTimeoutObj = setTimeout(function() {
//             //   socket.close();
//               // createWebSocket();
//             // }, self.timeout);
//         }, this.timeout)
//     },
// }

const heartCheck = {
    timeout: 5000,
    timeoutObj: setTimeout(() => { }),
    serverTimeoutObj: setInterval(() => { }),
    reset: function () {
        clearTimeout(this.timeoutObj);
        clearTimeout(this.serverTimeoutObj);
        return this;
    },
    start: function () {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const self = this;
        clearTimeout(this.timeoutObj);
        clearTimeout(this.serverTimeoutObj);
        this.timeoutObj = setTimeout(function () {
            socket.send(
                JSON.stringify({
                    event: 'ping',
                })
            );
            // self.serverTimeoutObj = setTimeout(function () {
            //     console.log('关闭服务');
            //     socket.close();
            // }, self.timeout);
        }, this.timeout);
    },
}
