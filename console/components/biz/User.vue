<template>
  <div class="card__body fn__clear">
    <v-form ref="form" @submit.prevent="created">
      <v-text-field
        :label="$t('account', $store.state.locale)"
        v-model="name"
        :counter="32"
        :rules="requiredRules"
        required
      ></v-text-field>
      <!-- 错误消息显示区域 -->
      <div class="alert alert--danger" v-show="error">
        <v-icon>danger</v-icon>
        <span>{{ errorMsg }}</span>
      </div>
    </v-form>
    <v-btn class="fn__right btn--margin-t30 btn--info btn--space" @click="created">
      {{ $t('confirm', $store.state.locale) }}
    </v-btn>
    <!-- 关闭此控件。搭配父组件的 :show.sync -->
    <v-btn class="fn__right btn--margin-t30 btn--danger btn--space" @click="$emit('update:show', false)">
      {{ $t('cancel', $store.state.locale) }}
    </v-btn>
  </div>
</template>

<script>
  import {required, maxSize} from '~/plugins/validate'

  export default {
    data () {
      return {
        errorMsg: '', // 错误消息
        error: false, // 是否出错
        name: '', // 用户名
        requiredRules: [
          (v) => required.call(this, v), // 必填
          (v) => maxSize.call(this, v, 32) // 长度最大 32 个字符
        ]
      }
    },
    methods: {
      // 创建
      async created () {
        if (!this.$refs.form.validate()) {
          return
        }
        const responseData = await this.axios.post('/console/users', {
          name: this.name
        })

        if (responseData.code === 0) {  // 成功
          this.$set(this, 'error', false)
          this.$set(this, 'errorMsg', '')
          this.$emit('addSuccess')
        } else { // 出错
          this.$set(this, 'error', true)
          this.$set(this, 'errorMsg', responseData.msg)
        }
      }
    }
  }
</script>
