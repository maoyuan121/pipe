import vueAxios from '~/plugins/axios'

export const state = () => ({
  locale: 'zh_CN', // 选择的语言
  version: '1.0.0', // 版本号
  isInit: false, // 是否初始
  name: '',
  nickname: '',
  blogTitle: '',
  avatarURL: '',
  blogURL: '/',
  role: 0, // 0-no login, 1-admin, 2-blog admin, 3-blog user, 4-visitor
  blogs: [{
    title: '',
    id: ''
  }],
  snackMsg: '', // 弹出层的消息内容
  snackBar: false, // 是否显示弹出层
  snackModify: 'error', // 弹出层的类型：成功、失败
  menu: [], // 管理页面的左侧菜单
  tagsItems: [],
  bodySide: ''
})

// mutation 只能有同步方法
// 用于改变 store 的 state
export const mutations = {
  setMenu (state, data) {
    state.menu = data
  },
  setBodySide (state, data) {
    state.bodySide = data
  },
  setLogout (state, data) {
    state.role = data
  },
  setStatus (state, data) {
    state.locale = data.locale
    state.version = data.version
    state.isInit = data.inited
    state.role = data.role
    state.name = data.name
    state.nickname = data.nickname
    state.blogTitle = data.blogTitle
    state.blogURL = data.blogURL
    state.blogs = data.blogs
    state.avatarURL = data.avatarURL
  },
  setLocale (state, locale) {
    state.locale = locale
  },
  setIsInit (state, isInit) {
    state.isInit = isInit
  },
  setBlog (state, data) {
    state.blogTitle = data.title
    state.blogURL = data.path
    state.role = data.role
  },
  setSnackBar (state, data) {
    state.snackBar = data.snackBar
    state.snackMsg = data.snackMsg
    if (data.snackModify) {
      state.snackModify = data.snackModify
    } else {
      state.snackModify = 'error'
    }
  },
  setTagsItems (state, data) {
    state.tagsItems = data
  }
}

// action 可以有异步方法
// action 提交的是 mutation，而不是直接变更状态
export const actions = {
  async nuxtClientInit ({ commit, state }, { app }) {
    // TrimB3Id
    const search = location.search
    if (search.indexOf('b3id') > -1) {
      history.replaceState('', '', window.location.href.replace(/(&b3id=\w{8})|(b3id=\w{8}&)|(\?b3id=\w{8}$)/, ''))
    }

    try {
      const responseData = await vueAxios().get('/status?' + (new Date()).getTime())
      if (responseData) {
        if (app.i18n.messages[responseData.locale]) {
          app.i18n.locale = responseData.locale
        } else {
          const message = require(`../../i18n/${responseData.locale}.json`)
          app.i18n.setLocaleMessage(responseData.locale, message)
        }
        commit('setStatus', responseData)
      } else {
        const message = require(`../../i18n/${state.locale}.json`)
        app.i18n.setLocaleMessage(state.locale, message)
      }
    } catch (e) {
      console.error(e)
      const message = require(`../../i18n/${state.locale}.json`)
      app.i18n.setLocaleMessage(state.locale, message)
    }
  },
  setLocaleMessage ({ commit }, locale) {
    if (this.app.i18n.messages[locale]) {
      this.app.i18n.locale = locale
    } else {
      const message = require(`../../i18n/${locale}.json`)
      this.app.i18n.setLocaleMessage(locale, message)
    }
    commit('setLocale', locale)
  },
  async getTags ({ commit, state }) {
    if (state.tagsItems.length > 0) {
      return
    }
    const tagResponseData = await vueAxios().get('/console/tags/')
    if (tagResponseData) {
      let tagList = []
      tagResponseData.map((v) => {
        tagList.push(v.title)
      })
      commit('setTagsItems', tagList)
    }
  }
}
