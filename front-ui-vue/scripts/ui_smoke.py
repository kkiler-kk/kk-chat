"""kk-chat 前端 UI 冒烟：登录页渲染 → 登录 → 会话/聊天 → 发消息 → 联系人面板"""
import subprocess
import sys

from playwright.sync_api import sync_playwright

results = []


def check(name, ok, extra=""):
    results.append(ok)
    print(f"{'✅' if ok else '❌'} {name} {extra}")


def redis_cmd(*args):
    out = subprocess.run(
        ["docker", "exec", "blog-redis", "redis-cli", "-n", "2", *args],
        capture_output=True, text=True,
    )
    return out.stdout.strip()


def get_captcha_answer():
    keys = [k for k in redis_cmd("KEYS", "captcha:*").splitlines() if k]
    if not keys:
        return None
    return redis_cmd("GET", keys[-1]).strip('"')


with sync_playwright() as p:
    browser = p.chromium.launch(headless=True)
    page = browser.new_page(viewport={"width": 1440, "height": 900})
    errors = []
    page.on("console", lambda m: errors.append(m.text) if m.type == "error" else None)
    page.on("pageerror", lambda e: errors.append(str(e)))

    # 1. 登录页渲染
    page.goto("http://localhost:9234/login")
    page.wait_for_load_state("networkidle")
    page.screenshot(path="/tmp/ui-1-login.png")
    check("登录页品牌渲染", page.locator(".brand").inner_text() == "kk-chat")
    check("图形验证码已加载", page.locator(".captcha-img").count() == 1)

    # 2. 登录 kk002
    answer = get_captcha_answer()
    check("从 redis 取得验证码答案", answer is not None, f"answer={answer}")
    page.fill('input[placeholder="identity 或邮箱"]', "kk002")
    page.fill('input[placeholder="密码"]', "123456")
    page.fill('input[placeholder="验证码"]', answer or "0000")
    page.click('.login-page button[type=submit]')
    page.wait_for_url("**/home", timeout=8000)
    page.wait_for_load_state("networkidle")
    page.wait_for_timeout(1500)
    page.screenshot(path="/tmp/ui-2-home.png")
    check("登录成功进入 /home", "/home" in page.url)

    # 3. 会话列表
    conv_names = page.locator(".conv-name").all_inner_texts()
    check("会话列表渲染", len(conv_names) > 0, f"convs={conv_names}")
    check("好友 KK 在会话列表", any("KK" in n for n in conv_names))

    # 4. 打开会话，历史消息渲染
    page.locator(".conv-item", has=page.locator(".conv-name", has_text="KK")).first.click()
    page.wait_for_timeout(1200)
    bubbles = page.locator(".bubble").count()
    check("聊天窗口历史消息渲染", bubbles > 0, f"bubbles={bubbles}")
    page.screenshot(path="/tmp/ui-3-chat.png")

    # 5. 发送消息
    page.fill('textarea[placeholder*="输入消息"]', "UI 冒烟测试消息")
    page.click('.send-row button')
    page.wait_for_timeout(1000)
    check("发送的消息上屏", page.locator(".bubble", has_text="UI 冒烟测试消息").count() >= 1)
    page.screenshot(path="/tmp/ui-4-sent.png")

    # 6. 联系人面板
    page.click('.nav-item[title="联系人"]')
    page.wait_for_timeout(800)
    friends = page.locator(".friend-row .fname").all_inner_texts()
    check("联系人面板好友列表", len(friends) > 0, f"friends={friends}")
    page.screenshot(path="/tmp/ui-5-contacts.png")

    # 7. 群组面板
    page.click('.nav-item[title="群组"]')
    page.wait_for_timeout(800)
    groups = page.locator(".group-row .gname").all_inner_texts()
    check("群组面板列表", len(groups) > 0, f"groups={groups}")
    page.screenshot(path="/tmp/ui-6-groups.png")

    # 8. 设置面板
    page.click('.nav-item[title="设置"]')
    page.wait_for_timeout(500)
    check("设置面板渲染", page.locator(".id-account").count() == 1)
    page.screenshot(path="/tmp/ui-7-settings.png")

    js_errors = [e for e in errors if "favicon" not in e.lower()]
    check("无 JS 运行时错误", len(js_errors) == 0, str(js_errors[:3]))

    browser.close()

print(f"\n=== UI 冒烟结果: {sum(results)}/{len(results)} 通过 ===")
sys.exit(0 if all(results) else 1)
