<template>
  <nav class="navbar navbar-expand-lg">
    <div class="container-fluid">
      <a class="navbar-brand" href="#">{{ $t('routes.' + routeName) }}</a>
      <button type="button"
              class="navbar-toggler navbar-toggler-right"
              :class="{toggled: $sidebar.showSidebar}"
              aria-controls="navigation-index"
              aria-expanded="false"
              aria-label="Toggle navigation"
              @click="toggleSidebar">
        <span class="navbar-toggler-bar burger-lines"></span>
        <span class="navbar-toggler-bar burger-lines"></span>
        <span class="navbar-toggler-bar burger-lines"></span>
      </button>
      <div class="collapse navbar-collapse justify-content-end">
        <ul class="navbar-nav ml-auto">
          <li class="nav-item" v-if="this.uuid">
            <router-link class="nav-link" :to="{ name: 'MemberEdit', params: { uuid: uuid}}">
              {{ $t('general.profile_zone') }}
            </router-link>
          </li>
          <li class="nav-item" v-if="!this.uuid">
            <router-link class="nav-link" :to="{ name: 'Login'}">
              {{ $t('general.login') }}
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<i18n src='assets/translations/routes.json'></i18n>
<i18n src='assets/translations/general.json'></i18n>

<script>
  import {mapGetters} from 'vuex'

  export default {
    computed: {
      ...mapGetters(['uuid', 'code', 'type']),
      routeName () {
        const {path} = this.$route
        return path.split('/')[1].toLowerCase()
      }
    },
    data () {
      return {
        activeNotifications: false
      }
    },
    methods: {
      capitalizeFirstLetter (string) {
        return string.charAt(0).toUpperCase() + string.slice(1)
      },
      toggleNotificationDropDown () {
        this.activeNotifications = !this.activeNotifications
      },
      closeDropDown () {
        this.activeNotifications = false
      },
      toggleSidebar () {
        this.$sidebar.displaySidebar(!this.$sidebar.showSidebar)
      },
      hideSidebar () {
        this.$sidebar.displaySidebar(false)
      },
      editMemberUuid (memberUuid) {
        this.$router.push({path: `/MemberEdit/${memberUuid}`})
      }
    }
  }

</script>
<style>

</style>
