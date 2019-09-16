<template>
  <div class="card__body fn__clear">
    <v-form ref="form">
      <v-text-field
        :label="$t('title', $store.state.locale)"
        v-model="title"
        :counter="128"
        :rules="titleRules"
        required
      ></v-text-field>
      <v-text-field
        label="URI"
        v-model="url"
        :rules="URIRules"
        :counter="255"
        required
      ></v-text-field>
      <v-text-field
        :label="$t('description', $store.state.locale)"
        v-model="description"
        :rules="descriptionRules"
        :counter="255"
      ></v-text-field>
      <v-select
        v-model="tags"
        :label="$t('tags', $store.state.locale)"
        chips
        tags
        :items="$store.state.tagsItems"
        required
        :rules="tagsRules"
      ></v-select>
      <div class="alert alert--danger" v-show="error">
        <v-icon>danger</v-icon>
        <span>{{ errorMsg }}</span>
      </div>
    </v-form>
    <v-btn class="fn__right btn--margin-t30 btn--info btn--space" @click="created">
      {{ $t('confirm', $store.state.locale) }}
    </v-btn>
    <v-btn class="fn__right btn--margin-t30 btn--danger btn--space" @click="$emit('update:show', false)">
      {{ $t('cancel', $store.state.locale) }}
    </v-btn>
  </div>
</template>

<script>
  import { required, maxSize } from '~/plugins/validate'

  export default {
    props: {
      id: {
        type: Number,
        required: true
      }
    },
    data () {
      return {
        errorMsg: '', // 错误消息
        error: false, // 是否发生错误
        title: '', // 类别名
        url: '', // url
        description: '', // 描述
        tags: '', // 标签
        titleRules: [ // 验证标题
          (v) => required.call(this, v), // 必填
          (v) => maxSize.call(this, v, 128) // 长度不得超过 128 个字符
        ],
        descriptionRules: [ // 验证描述
          (v) => maxSize.call(this, v, 255) // 长度不得超过 255 个字符
        ],
        URIRules: [ // 验证 URL
          (v) => required.call(this, v), // 必填
          (v) => maxSize.call(this, v, 255) // 长度不得超过 255 个字符
        ],
        tagsRules: [
          (v) => this.tags.length > 0 || this.$t('required', this.$store.state.locale)
        ]
      }
    },
    watch: {
      id: function () {
        this.init()
      }
    },
    methods: {
      // 创建/编辑 类别
      async created () {
        if (!this.$refs.form.validate()) {
          return
        }
        let responseData = {}
        const requestData = {
          title: this.title,
          path: this.url,
          description: this.description,
          tags: this.tags.join(',')
        }
        if (this.id === 0) {
          responseData = await this.axios.post('/console/categories', requestData)
        } else {
          responseData = await this.axios.put(`/console/categories/${this.id}`, requestData)
        }

        if (responseData.code === 0) {
          this.$set(this, 'error', false)
          this.$set(this, 'errorMsg', '')
          this.$emit('addSuccess') // 触发定义在该组建上的 addSuccess 事件。 <category v-if="showForm" :show.sync="showForm" @addSuccess="addSuccess" :id="editId"></category>
        } else {
          this.$set(this, 'error', true)
          this.$set(this, 'errorMsg', responseData.msg)
        }
      },
      // 如果 id > 0 从数据库里面取出这个类别的信息赋给相应的 data 属性
      async init () { 
        if (this.id === 0) {
          return
        }
        const responseData = await this.axios.get(`/console/categories/${this.id}`)
        if (responseData) {
          this.$set(this, 'title', responseData.title)
          this.$set(this, 'url', responseData.path)
          this.$set(this, 'description', responseData.description)
          this.$set(this, 'tags', responseData.tags.split(','))
        }
      }
    },
    mounted () {
      this.init() 
      // get tags
      // 异步操作，例如向后台提交数据，写法： this.$store.dispatch('actionName',值)
      this.$store.dispatch('getTags')
    }
  }
</script>
