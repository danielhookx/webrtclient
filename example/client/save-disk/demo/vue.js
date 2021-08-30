const comp = {
    data() {
        return {
            name: 'save-disk',
            server: '172.16.101.131:19801'
        }
    }
}

const vm = Vue.createApp(comp).mount('#app')