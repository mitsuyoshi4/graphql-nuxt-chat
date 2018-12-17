export const state = () => ({
  user: ''
})

export const mutations = {
  join(state, text) {
    state.user = text
  }
}
