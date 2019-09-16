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
        :label="$t('links', $store.state.locale)"
        v-model="url"
        :rules="linkRules"
        :counter="255"
        required
      ></v-text-field>
      <v-text-field
        :label="$t('iconPath', $store.state.locale)"
        v-model="iconURL"
        :rules="iconURLRules"
        :counter="255"
      ></v-text-field>
      <v-select
        :label="$t('openMethod', $store.state.locale)"
        :items="openMethods"
        v-model="openMethod"
        append-icon=""
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
        errorMsg: '', // 错误信息
        error: false, // 是否出错
        title: '', // 导航名
        url: '', // url
        iconURL: '', // 图标 url
        openMethod: '', // 打开方式。原地打开、新开、父跳转、top 跳转
        openMethods: [
          {
            'text': this.$t('openMethod1', this.$store.state.locale),
            'value': ''
          },
          {
            'text': this.$t('openMethod2', this.$store.state.locale),
            'value': '_blank'
          },
          {
            'text': this.$t('openMethod3', this.$store.state.locale),
            'value': '_parent'
          },
          {
            'text': this.$t('openMethod4', this.$store.state.locale),
            'value': '_top'
          }
        ],
        // 导航名的验证规则
        titleRules: [
          (v) => required.call(this, v), // 必填
          (v) => maxSize.call(this, v, 128) // 长度最多 128 个字符
        ],
        // 连接地址的验证规则
        linkRules: [
          (v) => required.call(this, v), // 必填
          (v) => maxSize.call(this, v, 255) // 长度最多 255 个字符
        ],
        // 图标的验证规则
        iconURLRules: [
          (v) => maxSize.call(this, v, 255) // 长度最多 255 个字符
        ]
      }
    },
    watch: {
      id: function () {
        this.init()
      }
    },
    methods: {
      // 新建|编辑导航
      async created () {
        if (!this.$refs.form.validate()) {
          return
        }
        let responseData = {}
        const requestData = {
          title: this.title,
          url: this.url,
          iconURL: this.iconURL,
          openMethod: this.openMethod
        }
        if (this.id === 0) {
          responseData = await this.axios.post('/console/navigations', requestData)
        } else {
          responseData = await this.axios.put(`/console/navigations/${this.id}`, requestData)
        }

        if (responseData.code === 0) {
          this.$set(this, 'error', false)
          this.$set(this, 'errorMsg', '')
          this.$emit('addSuccess')
        } else {
          this.$set(this, 'error', true)
          this.$set(this, 'errorMsg', responseData.msg)
        }
      },
      // 根据 id 获取导航
      async init () {
        if (this.id === 0) {
          return
        }
        const responseData = await this.axios.get(`/console/navigations/${this.id}`)
        if (responseData) {
          this.$set(this, 'title', responseData.title)
          this.$set(this, 'url', responseData.url)
          this.$set(this, 'iconURL', responseData.iconURL)
          this.$set(this, 'openMethod', responseData.openMethod)
        }
      }
    },
    mounted () {
      this.init()
    }
  }
</script>
