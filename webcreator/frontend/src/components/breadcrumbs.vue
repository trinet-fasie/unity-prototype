<template>
  <ol class="breadcrumb">
    <li
      v-for="(routeObject, index) in routeRecords"
      :key="index"
      class="breadcrumb-item">
      <span
        v-if="isLast(index)"
        class="active">{{ $t(routeObject.meta.label) }}</span>
      <router-link
        v-else
        :to="routeObject">{{ $t(routeObject.meta.label) }}</router-link>
    </li>
  </ol>
</template>

<script>
export default {
  props: {
    list: {
      type: Array,
      required: true,
      default: () => []
    }
  },
  computed: {
    routeRecords: function () {
      return this.list
    }
  },
  methods: {
    isLast (index) {
      return index === this.list.length - 1
    }
  }
}
</script>
