<template>
  <div class="modal-card">
    <section class="modal-card-body">
      <h3 class="title has-text-centered has-text-dark">Create User</h3>
      <div class="box">
        <b-field label="UserName">
          <b-input v-model="user" type="user" placeholder="user name">
          </b-input>
        </b-field>

        <button class="button is-dark is-large is-fullwidth" @click="join()">
          Join
        </button>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: 'LoginModal',
  data () {
    return {
      user: '',
    }
  },
  methods: {
    join() {
      this.$store.commit('users/join', this.user)
      this.$router.push('/chat')
    }
  }
}
</script>
