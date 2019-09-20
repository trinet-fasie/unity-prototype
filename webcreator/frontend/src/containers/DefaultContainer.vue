<template>
  <div class="app">
    <AppHeader fixed>
      <SidebarToggler
        class="d-lg-none"
        display="md"
        mobile />
      <b-link
        class="navbar-brand"
        to="/">
      </b-link>
      <SidebarToggler
        class="d-md-down-none"
        display="lg" />
      <b-navbar-nav class="ml-auto mr-4">
        <locale-changer />
      </b-navbar-nav>
    </AppHeader>
    <div class="app-body">
      <AppSidebar fixed>
        <SidebarHeader />
        <SidebarForm />
        <SidebarNav :nav-items="nav"/>
        <SidebarFooter />
        <SidebarMinimizer />
      </AppSidebar>
      <main class="main">
        <breadcrumbs :list="list" />
        <div class="container-fluid">
          <router-view/>
        </div>
      </main>
      <AppAside fixed>
        <!--aside-->
        <DefaultAside />
      </AppAside>
    </div>
    <TheFooter>
    </TheFooter>
  </div>
</template>

<script>
import nav from '@/_nav'
import {
  Aside as AppAside,
  AsideToggler,
  Breadcrumb,
  Footer as TheFooter,
  Header as AppHeader,
  Sidebar as AppSidebar,
  SidebarFooter,
  SidebarForm,
  SidebarHeader,
  SidebarMinimizer,
  SidebarNav,
  SidebarToggler
} from '@coreui/vue'
import DefaultAside from './DefaultAside'
import BLink from 'bootstrap-vue/es/components/link/link'
import BNavbarNav from 'bootstrap-vue/es/components/navbar/navbar-nav'
import BNavItem from 'bootstrap-vue/es/components/nav/nav-item'
import BBadge from 'bootstrap-vue/es/components/badge/badge'
import Breadcrumbs from '../components/breadcrumbs'
import LocaleChanger from '../components/locale-changer'

export default {
  components: {
    BLink,
    BNavbarNav,
    BNavItem,
    BBadge,
    AsideToggler,
    AppHeader,
    AppSidebar,
    AppAside,
    TheFooter,
    Breadcrumb,
    DefaultAside,
    SidebarForm,
    SidebarFooter,
    SidebarToggler,
    SidebarHeader,
    SidebarNav,
    Breadcrumbs,
    LocaleChanger,
    SidebarMinimizer
  },
  data () {
    return {
      nav: nav.items
    }
  },
  computed: {
    name () {
      return this.$route.name
    },
    list () {
      return this.$route.matched.filter((route) => route.meta.label)
    },
    currentYear () {
      return new Date().getFullYear()
    }
  }
}
</script>
