<template>
  <div class="card">
    <user v-if="showForm" :show.sync="showForm" @addSuccess="addSuccess"></user>

    <div v-show="!showForm" class="card__body fn__flex">
      <v-text-field
        v-if="list.length > 0"
        @keyup.enter="getList()"
        class="fn__flex-1"
        :label="$t('enterSearch', $store.state.locale)"
        v-model="keyword">
      </v-text-field>
      <v-btn class="btn--success" :class="{'btn--new': list.length > 0}" @click="edit">{{ $t('new', $store.state.locale) }}</v-btn>
    </div>

    <ul class="list" v-if="list.length > 0">
      <li v-for="item in list" :key="item.id" class="fn__flex"
          v-if="($store.state.role === 3 && item.name === $store.state.name) || $store.state.role < 3">
        <a :href="item.url"
           :aria-label="item.name"
           class="avatar avatar--mid avatar--space pipe-tooltipped pipe-tooltipped--n"
           :style="`background-image: url(${item.avatarURL})`"></a>
        <div class="fn__flex-1">
          <div class="fn__flex">
            <a class="list__title fn__flex-1" :href="item.url">
              {{ item.nickname || item.name }}
            </a>
            <v-btn class="btn--small btn--info" @click="prohibit(item.id, 'unprohibit')" v-if="item.role === 4">
              {{ $t('unProhibit', $store.state.locale) }}
            </v-btn>
            <v-btn class="btn--small btn--danger" @click="prohibit(item.id, 'prohibit')" v-else>
              {{ $t('prohibit', $store.state.locale) }}
            </v-btn>
          </div>
          <div class="list__meta">
            <span class="fn-nowrap">{{ item.articleCount }} {{ $t('article', $store.state.locale) }}</span> •
            <span class="fn-nowrap" :class="{'ft__danger': item.role === 4}">{{ getRoleName(item.role) }}</span>
          </div>
        </div>
      </li>
    </ul>
    <div class="pagination--wrapper fn__clear" v-if="pageCount > 1">
      <v-pagination
        :length="pageCount"
        v-model="currentPageNum"
        :total-visible="windowSize"
        class="fn__right"
        circle
        next-icon="angle-right"
        prev-icon="angle-left"
        @input="getList"
      ></v-pagination>
    </div>
  </div>
</template>

<script>
  import User from '~/components/biz/User'

  export default {
    components: {
      User
    },
    data () {
      return {
        showForm: false, // 是否显示表单
        currentPageNum: 1, // 当前页码
        pageCount: 1, // 总页数
        windowSize: 1, // 显示多少个分页按钮
        list: [], // 用户列表
        keyword: '' // 搜索关键字
      }
    },
    head () {
      return {
        title: `${this.$t('userList', this.$store.state.locale)} - ${this.$store.state.blogTitle}`
      }
    },
    methods: {
      // 获取角色名
      // 1：超级管理员，2:博客管理员，3:博客用户，4：被禁用的用户，默认为博客用户
      getRoleName (role) {
        let roleName = this.$t('blogUser', this.$store.state.locale)
        switch (role) {
          case 1:
            roleName = this.$t('superAdmin', this.$store.state.locale)
            break
          case 2:
            roleName = this.$t('blogAdmin', this.$store.state.locale)
            break
          case 3:
            roleName = this.$t('blogUser', this.$store.state.locale)
            break
          case 4:
            roleName = this.$t('prohibitUser', this.$store.state.locale)
            break
          default:
            break
        }
        return roleName
      },
      // 获取用户列表
      async getList (currentPage = 1) {
        const responseData = await this.axios.get(`/console/users?p=${currentPage}&key=${this.keyword}`)
        if (responseData) {
          this.$set(this, 'list', responseData.users)
          this.$set(this, 'currentPageNum', responseData.pagination.currentPageNum)
          this.$set(this, 'pageCount', responseData.pagination.pageCount)
          this.$set(this, 'windowSize', document.documentElement.clientWidth < 721 ? 5 : responseData.pagination.windowSize)
        }
      },
      // 禁用|解禁
      // type 为 prohibit 等于禁用，为 unprohibit 等于解禁
      async prohibit (id, type) {
        const responseData = await this.axios.put(`/console/users/${id}/${type}`)
        if (responseData.code === 0) {
          this.$set(this, 'error', false)
          this.$set(this, 'errorMsg', '')
          this.getList()
        } else {
          this.$set(this, 'error', true)
          this.$set(this, 'errorMsg', responseData.msg)
        }
      },
      // 添加成功事件的 handler
      // 重新获取用户集合，隐藏表单
      addSuccess () {
        this.getList()
        this.$set(this, 'showForm', false)
      },
      // 显示表单
      edit () {
        this.$set(this, 'showForm', true)
      }
    },
    mounted () {
      this.getList()
    }
  }
</script>
