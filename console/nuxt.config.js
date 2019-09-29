import Vuetify from "vuetify";
import VueI18n from "./plugins/init";

const env = require(`../pipe.json`)

module.exports = {
  /*
  ** 该配置项用于定义应用客户端和服务端的环境变量。
  ** 详见：https://zh.nuxtjs.org/api/configuration-env
  */
  env,
  /*
  ** 该配置项用于配置应用默认的meta标签。
  ** 详见：https://zh.nuxtjs.org/api/configuration-head
  */
  head: {
    title: 'Pipe',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { name: 'robots', content: 'none' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      { rel: 'manifest', href: `/manifest.json` }
    ]
  },
  /*
  ** 该配置项用于个性化定制 Nuxt.js 使用的加载组件。
  ** 详见：https://zh.nuxtjs.org/api/configuration-loading
  */
  loading: { color: '#4a4a4a' },
  /*
  ** 该配置项用于定义应用的全局（所有页面均需引用的）样式文件、模块或第三方库。
  ** 详见：https://zh.nuxtjs.org/api/configuration-css/
  */
  css: [
    'vuetify/dist/vuetify.min.css',
    '~assets/scss/main.scss'
  ],
  /*
  ** 该配置项用于配置的所有插件会在 Nuxt.js 应用初始化之前被加载导入
  ** ssr: Boolean (默认为 true) 如果值为 false，该文件只会在客户端被打包引入。
  ** https://zh.nuxtjs.org/api/configuration-plugins
  */
  plugins: [
    { src: '~/plugins/axios.js', ssr: false }, // 个性化下 axios， eg：代理地址、请求出错的时候弹出错误提示信息层
    { src: '~/plugins/init.js', ssr: false }, // 使用 Vuetify VueI18n
    { src: '~/plugins/nuxt-client-init.js', ssr: false } // 加载多语言，获取当前用户等信息
  ],
  mode: 'spa',
  /*
  ** Build configuration
  ** 详见：https://zh.nuxtjs.org/api/configuration-build
  */
  build: {
    publicPath: (env.StaticServer ||  env.Server) + '/console/dist/',
    extractCSS: true,
    ssr: false,
    /*
    ** Run ESLint on save
    */
    extend (config, ctx) {
      if (ctx.dev && ctx.isClient) {
        config.module.rules.push({
          enforce: 'pre',
          test: /\.(js|vue)$/,
          loader: 'eslint-loader',
          exclude: /(node_modules)/
        })
      }
    }
  },
  /*
  ** 该配置项可用于覆盖 Nuxt.js 默认的 vue-router 配置。
  ** 详见：https://zh.nuxtjs.org/api/configuration-router
  */
  // router: {
  //   middleware: ['authenticated']
  // },
  /*
  ** 该配置项允许您将Nuxt模块添加到项目中。
  ** 详见：https://zh.nuxtjs.org/api/configuration-modules
  */
  modules: ['@nuxtjs/proxy'],
  proxy: {
    '/api': {
      target: env.Server,
      changeOrigin: true
    },
    '/mock': {
      target: env.MockServer,
      changeOrigin: true,
      pathRewrite: {
        '^/mock/': ''
      }
    }
  }
}
