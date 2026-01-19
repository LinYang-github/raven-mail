import { registerMicroApps, start } from 'qiankun';

// 1. 定义子应用
const apps = [
  {
    name: 'raven-mail', // 子应用名称
    entry: '//localhost:5173', // 子应用开发环境地址 (Raven web dev server)
    container: '#subapp-container', // 挂载容器 ID
    activeRule: '/mail', // 激活路由
    props: {
        // 全局传值示例：模拟 Token 下发
        token: 'demo-host-token-xyz',
        user: { name: 'Host Admin', id: 'admin-001' }
    }
  },
];

// 2. 注册子应用
registerMicroApps(apps);

// 3. 启动
start({
    prefetch: false, // Demo环境通常关闭预加载
    sandbox: {
        experimentalStyleIsolation: true // Use experimental scope isolation instead
    }
});

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
