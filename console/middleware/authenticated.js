/*
 ** 如果没有初始化跳转到 start 进行初始化
 ** 如果已经初始化了但没有登录或者是访客那么跳转到
 *
*/

export default function ({redirect, store, route}) {
  // 如果应用还没有被初始化，并且 route 不是 /start，那么固定跳转到 /start
  if (!store.state.isInit) {
    if (route.path !== '/start') {
      return redirect('/start')
    }
    return
  }

  const isLogin = store.state.role !== 0
  if (route.path.indexOf('/admin') > -1) { // 如果 router 是 /admin，并且没有登录或者是游客，那么跳转到 /
    if (!isLogin || store.state.role === 4) {
      redirect('/')
    }
  } else if (route.path === '/start') { // 如果路径是 /start 并且已经初始化过了，并且已经登录了，那么跳转到 /admin 去
    if (isLogin) {
      redirect('/admin')
    }
  }
}
