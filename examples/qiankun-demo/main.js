import { registerMicroApps, start, initGlobalState } from 'qiankun';

// 初始状态
const state = {
    user: { name: 'User A', id: 'user-123', role: 'RED' }
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

// 监听状态变更并更新局部 state
actions.onGlobalStateChange((newState) => {
    console.log('[host] state changed', newState);
    state.user = newState.user;
});

// 1. 定义子应用
const apps = [
  {
    name: 'raven-mail', // 子应用名称
    entry: '//localhost:5173', // 子应用开发环境地址
    container: '#subapp-container', // 挂载容器 ID
    activeRule: '/mail', // 激活路由
    props: {
        token: 'demo-host-token-xyz',
        user: state.user,
        // 关键点：将当前的 role 闭包进去或通过 state 实时获取
        fetchUsers: (query) => fetchUsersMockByRole(query, state.user)
    }
  },
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

// User Switch Logic
const userMap = {
    'user-123': { name: 'User A', role: 'RED' },
    'user-456': { name: 'User B', role: 'BLUE' },
    'admin': { name: 'Admin', role: 'WHITE' }
};

const selector = document.getElementById('user-select');
if (selector) {
    selector.addEventListener('change', (e) => {
        const uid = e.target.value;
        const info = userMap[uid];
        const newUser = { id: uid, ...info };
        
        console.log('[host] switching user to:', newUser);
        actions.setGlobalState({ user: newUser });
    });
}


// Demo 简单的路由高亮逻辑
const updateNav = () => {
    const path = window.location.pathname;
    document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
    document.getElementById('home-content').style.display = path.startsWith('/mail') ? 'none' : 'block';
    
    if(path.startsWith('/mail')) {
        document.getElementById('nav-mail').classList.add('active');
    }
};

window.addEventListener('popstate', updateNav);
// 初始运行
updateNav();
