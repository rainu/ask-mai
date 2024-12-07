<template>
  <v-row dense class="pa-0 ma-0">
    <v-col :cols="progress ? 11 : 12" class="pa-0 ma-0">
      <v-text-field
        v-model="value"
        :disabled="progress"
        @keyup.enter="onSubmit"
        hide-details
        autofocus
        :placeholder="$t('prompt.placeholder')"
      >
      </v-text-field>
    </v-col>
    <v-col :cols="progress ? 1 : 0" v-show="progress" class="pa-0 ma-0">
      <v-btn @click="onStop" variant="flat" color="error" block style="height: 100%">
        <v-icon size="x-large">mdi-stop-circle</v-icon>
      </v-btn>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
	name: 'ChatInput',
  emits: ['update:modelValue', 'submit', 'interrupt'],
  props: {
    progress: {
      type: Boolean,
      required: false,
      default: false,
    },
    modelValue: {
      type: String,
      required: false,
      default: '',
    },
  },
  computed: {
    value: {
      get() {
        return this.modelValue
      },
      set(value: string) {
        this.$emit('update:modelValue', value)
      }
    }
  },
  methods: {
    onSubmit() {
      this.$emit('submit', this.modelValue)
    },
    onStop() {
      this.$emit('interrupt')
    },
  },
})
</script>

<style scoped></style>