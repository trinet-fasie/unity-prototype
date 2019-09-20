<template>
  <div class="mypad">
    <router-view/>
  </div>
</template>

<script>
import VBar from 'v-bar'

export default {
  components: { VBar }
}
</script>

<style>
  .mypad {
    width: 1024px;
    height: 1024px;
  }

  .mypad .container {
    padding: 15px 0 0 15px;
  }

  .wrapper {
    width: 1007px;
    height: 1024px;
    background: #eee;
  }

  .body {
    background: transparent;
  }
</style>
