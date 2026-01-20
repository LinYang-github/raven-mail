import { registerMicroApps, start, initGlobalState } from 'qiankun';

// 初始配置：支持通过 URL 参数或环境变量覆盖地址
const BACKEND_URL = window.localStorage.getItem('RAVEN_BACKEND_URL') || 'http://localhost:8080';
const SUBAPP_ENTRY = window.localStorage.getItem('RAVEN_SUBAPP_ENTRY') || '//localhost:5173';

const state = {
    user: { name: '红方-01', id: 'user-123', role: 'RED' },
    sessionId: 'default',
    themeColor: '#409EFF'
};

const allUsersInDirectory = [
    { id: 'user-123', name: '红方-01', dept: '红方阵营', role: 'RED' },
    { id: 'red-commander', name: '红方指挥官', dept: '红方阵营', role: 'RED' },
    { id: 'user-456', name: '蓝方-01', dept: '蓝方阵营', role: 'BLUE' },
    { id: 'blue-scout', name: '蓝方侦察兵', dept: '蓝方阵营', role: 'BLUE' },
    { id: 'admin', name: '总导演-01', dept: '导演部', role: 'WHITE' },
    { id: 'observer-01', name: '第三方观察员', dept: '白方部', role: 'WHITE' },
];

const fetchUsersMockByRole = async (query, currentUser) => {
    return new Promise(resolve => {
        setTimeout(() => {
            const filtered = allUsersInDirectory.filter(u => {
                // 1. 匹配搜索
                const isMatch = query === '' || u.name.includes(query) || u.id.includes(query);
                if (!isMatch) return false;

                // 2. 核心权限逻辑
                const myRole = currentUser.role || 'RED';
                if (myRole === 'WHITE') return true; // 白方看全员
                if (myRole === 'RED') return u.role !== 'BLUE'; // 红方不见蓝方
                if (myRole === 'BLUE') return u.role !== 'RED'; // 蓝方不见红方
                return true;
            });
            resolve(filtered);
        }, 300);
    });
};

// 初始化全局状态
const actions = initGlobalState(state);

// 提取 UI 更新函数
const updateUnreadBadge = (count) => {
    const badge = document.getElementById('mail-badge');
    if (badge) {
        badge.innerText = count;
        badge.style.display = count > 0 ? 'inline-block' : 'none';
    }
};

// 监听状态变更并更新局部 state
actions.onGlobalStateChange((newState) => {
    console.log('[host] state changed', newState);
    state.user = newState.user || state.user;
    state.sessionId = newState.sessionId || state.sessionId;
    state.themeColor = newState.themeColor || state.themeColor;
    
    // 处理未读数变更
    if (newState.unreadCount !== undefined) {
        updateUnreadBadge(newState.unreadCount);
    }
});

// 1. 定义子应用
// 1. 定义子应用 - 拆分为三个独立注册的实例以支持不同的 Base Route
const apps = [
  {
    name: 'raven-mail',
    entry: SUBAPP_ENTRY,
    container: '#subapp-container',
    activeRule: '/mail',
    props: {
        token: 'demo-host-token-xyz',
        getGlobalState: () => state,
        onGlobalStateChange: (callback, fireImmediately) => actions.onGlobalStateChange(callback, fireImmediately),
        ravenConfig: { showReset: true, showSidebar: true, primaryColor: state.themeColor },
        fetchUsers: (query) => fetchUsersMockByRole(query, state.user),
        
        // Feature Flags & Config
        modules: ['mail'],
        routeBase: '/mail'
    }
  },
  {
    name: 'raven-im',
    entry: SUBAPP_ENTRY,
    container: '#subapp-container',
    activeRule: '/im',
    props: {
        token: 'demo-host-token-xyz',
        getGlobalState: () => state,
        onGlobalStateChange: (callback, fireImmediately) => actions.onGlobalStateChange(callback, fireImmediately),
        ravenConfig: { showReset: true, showSidebar: true, primaryColor: state.themeColor },
        fetchUsers: (query) => fetchUsersMockByRole(query, state.user),
        
        // Feature Flags & Config
        modules: ['im'],
        routeBase: '/im'
    }
  },
  {
    name: 'raven-all',
    entry: SUBAPP_ENTRY,
    container: '#subapp-container',
    activeRule: '/all',
    props: {
        token: 'demo-host-token-xyz',
        getGlobalState: () => state,
        onGlobalStateChange: (callback, fireImmediately) => actions.onGlobalStateChange(callback, fireImmediately),
        ravenConfig: { showReset: true, showSidebar: true, primaryColor: state.themeColor },
        fetchUsers: (query) => fetchUsersMockByRole(query, state.user),
        
        // Feature Flags & Config
        modules: ['mail', 'im'],
        routeBase: '/all'
    }
  }
];

// 2. 注册子应用
registerMicroApps(apps);

// 3. 启动
start({
    prefetch: false,
    sandbox: {
        experimentalStyleIsolation: true
    }
});

// User Switch Logic - 保持与通讯录一致
const userMap = {
    'user-123': { name: '红方-01', role: 'RED' },
    'red-commander': { name: '红方指挥官', role: 'RED' },
    'user-456': { name: '蓝方-01', role: 'BLUE' },
    'blue-scout': { name: '蓝方侦察兵', role: 'BLUE' },
    'admin': { name: '总导演-01', role: 'WHITE' },
    'observer-01': { name: '第三方观察员', role: 'WHITE' }
};

const selector = document.getElementById('user-select');
if (selector) {
    // 重新生成选择器选项
    selector.innerHTML = Object.entries(userMap).map(([id, info]) => 
        `<option value="${id}">${info.name} (${info.role})</option>`
    ).join('');

    selector.addEventListener('change', (e) => {
        const uid = e.target.value;
        const info = userMap[uid];
        const newUser = { id: uid, ...info };
        
        console.log('[host] switching user to:', newUser);
        actions.setGlobalState({ user: newUser });
    });
}

const themeSelector = document.getElementById('theme-select');
if (themeSelector) {
    themeSelector.addEventListener('change', (e) => {
        const color = e.target.value;
        console.log('[host] switching theme color to:', color);
        actions.setGlobalState({ 
            themeColor: color,
            ravenConfig: {
                showReset: true,
                showSidebar: true,
                primaryColor: color
            }
        });
    });
}

const sessionSelector = document.getElementById('session-select');
if (sessionSelector) {
    sessionSelector.addEventListener('change', (e) => {
        const sid = e.target.value;
        console.log('[host] switching session to:', sid);
        actions.setGlobalState({ sessionId: sid });
    });
}

// 暴露重置场次接口
window.resetSession = async () => {
    const sid = state.sessionId;
    if (!sid) {
        alert('无有效场次 ID');
        return;
    }
    if (!confirm(`确定要重置场次 "${sid}" 吗？这将删除该场次下的所有邮件和文档数据。`)) {
        return;
    }

    try {
        const url = `${BACKEND_URL}/api/v1/sessions/${sid}`;
        const resp = await fetch(url, {
            method: 'DELETE'
        });
        
        const text = await resp.text();
        console.log('[host] Reset session response:', text);

        if (resp.ok) {
            alert('场次数据已清空');
            // 简单粗暴但有效：刷新页面重置所有子应用状态
            window.location.reload();
        } else {
            let errorMsg = text;
            try {
                const result = JSON.parse(text);
                errorMsg = result.error || text;
            } catch (e) {}
            alert(`清理失败: ${errorMsg}`);
        }
    } catch (err) {
        alert(`网络错误: ${err.message}`);
    }
};


// Demo 简单的路由高亮逻辑
const updateNav = () => {
    const path = window.location.pathname;
    document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
    document.getElementById('home-content').style.display = path.startsWith('/mail') ? 'none' : 'block';
    
    if(path.startsWith('/mail')) {
        document.getElementById('nav-mail').classList.add('active');
        // 进入邮件模块时不再强制清除徽标，保持其实时性
    }
};

// 监听子应用发送的新邮件自定义事件
window.addEventListener('raven-new-mail', (event) => {
    const { unreadCount, userId } = event.detail;
    console.log(`[host] New mail for ${userId}, unread: ${unreadCount}`);
    
    // 更新侧边栏徽标
    updateUnreadBadge(unreadCount);
    
    // 桌面通知示例
    if (Notification.permission === 'granted') {
        new Notification('新文电通知', {
            body: `您收到了一封新文电，请及时处理。`,
            icon: '/favicon.ico'
        });
    } else if (Notification.permission !== 'denied') {
        Notification.requestPermission();
    }
});

// 监听即时消息
window.addEventListener('raven-im-received', (event) => {
    const msg = event.detail;
    console.log(`[host] New IM from ${msg.sender_id}: ${msg.content}`);
    
    if (Notification.permission === 'granted') {
        new Notification('新即时消息', {
            body: `${msg.sender_id}: ${msg.content}`,
            icon: '/favicon.ico'
        });
    }
});

window.addEventListener('popstate', updateNav);
// 初始运行
updateNav();

