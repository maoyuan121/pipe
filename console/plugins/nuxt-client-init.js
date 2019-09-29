export default async (ctx) => {
  // dispatch：含有异步操作，数据提交至 actions ，可用于向后台提交数据
  // commit：同步操作，数据提交至 mutations ，可用于登录成功后读取用户信息写到缓存里
  await ctx.store.dispatch('nuxtClientInit', ctx)
}
