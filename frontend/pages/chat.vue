<template>
  <section class="section">
    <div>
      ユーザ一覧
      <b-taglist>
        <b-tag type="is-light" size="is-large" v-for="user in users" :key="user" :style="isSelf(user)">{{user}}</b-tag>
      </b-taglist>
    </div>
    <div>
      コメント
      <div v-for="message in messages" :key="message.id" class="message-wrapper is-clearfix">
        <b-collapse class="panel" open>
          <div slot="default" class="panel-heading" :style="isSelf(message.user)">
            <strong>{{ message.user }}</strong> <small>{{ message.createdAt }}</small>
          </div>
          <div class="panel-block">
            {{ message.text }}
          </div>
        </b-collapse>
      </div>
    </div>
    <hr>
    <div>
      <b-input type="textarea" v-model="message"></b-input>
      <p class="control">
        <button class="button is-primary" @click="addMessage">
          Send message
        </button>
      </p>
    </div>
  </section>
</template>

<script>
import QMessages from '@/apollo/queries/messages.gql'
import QUsers from '@/apollo/queries/users.gql'
import MMessagePost from '@/apollo/mutations/messagePost.gql'
import SMessagePosted from '@/apollo/subscriptions/messagePosted.gql'
import SUserjoined from '@/apollo/subscriptions/userJoined.gql'

export default {
  name: 'ChatPage',
  validate({params, query, store}) {
    return !!store.state.users.user
  },
  data() {
    return {
      message: ''
    }
  },
  methods: {
    addMessage: function() {
      this.$apollo.mutate({
        mutation: MMessagePost,
        variables: {
          user: this.$store.state.users.user,
          text: this.message,
        },
      })
      .then(() => {
        this.message = ''
      })
    },
    isSelf: function(user) {
      return (user === this.$store.state.users.user) ? { "background-color": "#00dddd" } : {}
    }
  },
  apollo: {
    messages: {
      query: QMessages,
      subscribeToMore: {
        document: SMessagePosted,
        variables () {
          return {
            user: this.$store.state.users.user
          }
        },
        updateQuery: (prev, { subscriptionData }) => {
          if (!subscriptionData.data) {
            return prev
          }
          const message = subscriptionData.data.messagePosted
          if (prev.messages.find(m => m.id === message.id)) {
            return prev
          }
          return Object.assign({}, prev, {
            messages: [message, ...prev.messages],
          })
        }
      },
    },
    users: {
      query: QUsers,
      subscribeToMore: {
        document: SUserjoined,
        variables () {
          console.log(this.$store.state.users.user)
          return {
            user: this.$store.state.users.user
          }
        },
        updateQuery: (prev, { subscriptionData }) => {
          if (!subscriptionData.data) {
            return prev
          }
          const user = subscriptionData.data.userJoined
          if (prev.users.find(u => u === user)) {
            return prev
          }
          return Object.assign({}, prev, {
            users: [user, ...prev.users],
          })
        }
      },
    }
  }
}
</script>

<style>
</style>
