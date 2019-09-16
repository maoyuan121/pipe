<template>
  <div class="card">
    <!-- 在子组件中 $emit('update:show', false) 等于将父的 shorForm 更新为了 false -->
    <category v-if="showForm" :show.sync="showForm" @addSuccess="addSuccess" :id="editId"></category>
    <!-- 如果没有显现上面的表单，那么就显示新建按钮 -->
    <div v-show="!showForm" class="card__body fn__clear">
      <v-btn class="btn--success" :class="{'fn__right': list.length > 0}" @click="edit(0)">{{ $t('new', $store.state.locale) }}</v-btn>
    </div>
    <ul class="list" v-if="list.length > 0">
      <li v-for="item in list" :key="item.id" class="fn__flex">
        <div class="fn__flex-1">
          <div class="fn__flex">
            <a class="list__title fn__flex-1"
               @click.stop="openURL(item.url)"
               href="javascript:void(0)">
              {{ item.title }}
            </a>
            <v-menu
              v-if="$store.state.role < 3"
              :nudge-bottom="28"
              :nudge-width="60"
              :nudge-left="60"
              :open-on-hover="true">
              <v-toolbar-title slot="activator">
                <v-btn class="btn--info btn--small" @click="edit(item.id)">
                  {{ $t('edit', $store.state.locale) }}
                  <v-icon>arrow_drop_down</v-icon>
                </v-btn>
              </v-toolbar-title>
              <v-list>
                <v-list-tile @click="edit(item.id)" class="list__tile--link">
                  {{ $t('edit', $store.state.locale) }}
                </v-list-tile>
                <v-list-tile @click="remove(item.id)" class="list__tile--link">
                  {{ $t('delete', $store.state.locale) }}
                </v-list-tile>
              </v-list>
            </v-menu>
          </div>
          <div class="list__meta">
            {{ item.description }}
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
  import Category from '~/components/biz/Category'

  export default {
    components: {
      Category
    },
    data () {
      return {
        editId: '', // 编辑的分类 id, 为 0 表示新建
        showForm: false, // 是否显示创建/编辑 Category 的表单
        currentPageNum: 1, // 当前页码 
        pageCount: 1, // 总页数
        windowSize: 1, // 设置可见页面按钮的最大数量
        list: [], // 分类集合
      }
    },
    head () {
      return {
        title: `${this.$t('categoryList', this.$store.state.locale)} - ${this.$store.state.blogTitle}`
      }
    },
    methods: {
      openURL (url) {
        window.location.href = url
      },
      // 获取分类集合
      async getList (currentPage = 1) {
        const responseData = await this.axios.get(`/console/categories?p=${currentPage}`)
        if (responseData) {
          // 赋值 data
          this.$set(this, 'list', responseData.categories || [])
          this.$set(this, 'currentPageNum', responseData.pagination.currentPageNum)
          this.$set(this, 'pageCount', responseData.pagination.pageCount)
          this.$set(this, 'windowSize', document.documentElement.clientWidth < 721 ? 5 : responseData.pagination.windowSize)
        }
      },
      // 删除分类
      async remove (id) {
        const responseData = await this.axios.delete(`/console/categories/${id}`)
        if (responseData === null) {
          // 弹出成功提示层
          this.$store.commit('setSnackBar', {
            snackBar: true,
            snackMsg: this.$t('deleteSuccess', this.$store.state.locale),
            snackModify: 'success'
          })
          // 重新加载分类集合
          this.getList()
          // 设置隐藏表单
          this.$set(this, 'showForm', false)
        }
      },
      // 创建/编辑成功后调用，子控件里面 $emit
      addSuccess () {
        this.getList()
        this.$set(this, 'showForm', false)
      },
      // 点击编辑
      edit (id) {
        // 显示表单
        this.$set(this, 'showForm', true)
        // 设置 editId，Category 子组建里面 watch 了 id, id 一变会从数据库里面获取指定的分类
        this.$set(this, 'editId', id)
      }
    },
    mounted () {
      // 加载分类集合
      this.getList()
    }
  }
</script>
