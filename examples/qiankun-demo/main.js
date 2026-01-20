import { registerMicroApps, start, initGlobalState } from 'qiankun';

// 初始状态
const state = {
    user: { name: 'User A', id: 'user-123' }
};

// 初始化全局状态
const actions = initGlobalState(state);

// 监听状态变更
actions.onGlobalStateChange((state, prev) => {
    console.log('[host] state changed', state);
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
        user: state.user, // Initial prop
        // 方案 1: 注入用户搜索函数 (模拟业务系统接口)
        fetchUsers: async (query) => {
            console.log('[host] fetchUsers called with:', query);
            const allUsers = [
                { id: 'user-123', name: 'User A', dept: '研发部' },
                { id: 'user-456', name: 'User B', dept: '设计部' },
                { id: 'admin', name: '系统管理员', dept: '管理层' },
                { id: 'zhangsan', name: '张三', dept: '财务部' },
                { id: 'lisi', name: '李四', dept: '市场部' },
            ];
            // 简单模拟搜索过滤
            return new Promise(resolve => {
                setTimeout(() => {
                    const filtered = query 
                        ? allUsers.filter(u => u.name.includes(query) || u.id.includes(query))
                        : allUsers; // 返回全部
                    resolve(filtered);
                }, 300); // 模拟网络延迟
            });
        }
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
    'user-123': 'User A',
    'user-456': 'User B',
    'admin': 'Admin'
};

const selector = document.getElementById('user-select');
if (selector) {
    selector.addEventListener('change', (e) => {
        const uid = e.target.value;
        const name = userMap[uid];
        const newUser = { name, id: uid };
        
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
