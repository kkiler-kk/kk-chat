import { Session } from "./storage";
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
        if (Session.get('token') == '') {
            throw new Error('用户尚未登录， 请稍后重试。。。')
        }
        socket = new WebSocket(import.meta.env.VITE_WEBSOCKET_URL + "?token=" + Session.get("token"))
        console.log('wsAdd', import.meta.env.VITE_WEBSOCKET_URL)
        init()
    } catch (e) {
        console.log('createSocket err:' + e);
        reconnect(); // 重连一下websocket
    }
    if (lockReconnect) {
        lockReconnect = false;
    }
};
export function init() {
    socket.onopen = function (_) {
        console.log('WebSocket:已连接');
        isActive = true;
    };

    socket.onmessage = function (event) {
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
    }

    socket.onerror = function (_) {
        console.log('WebSocket:发生错误');
        reconnect();
        isActive = false;
      };
  
      socket.onclose = function (_) {
        console.log('WebSocket:已关闭');
        reconnect();
        isActive = false;
      };
  
      window.onbeforeunload = function () {
        socket.close();
        isActive = false;
      };
};
const reconnect = () => {
    console.log('lockReconnect:' + lockReconnect);
    if (lockReconnect) return;
    lockReconnect = true;
    clearTimeout(timer);
    timer = setTimeout(() => {
        createSocket();
    }, 1000);
}
