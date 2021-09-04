<template>
  <div class="home">
    <Table :data="commands" :headers="this.headers"/>
  </div>
</template>

<script>
// @ is an alias to /src
import { mapState, mapActions } from 'vuex'
import Table from '../components/Table.vue';

export default {
  name: 'Commands',
  components: {
    Table
  },
  data: function() {
    return {headers: ['Command Name', 'Command Response', 'Command Perm', 'Command Count']}
  },
  computed: mapState({
    commands: state => state.commands
  }),
  methods: mapActions('commands', [
    'loadCommands'
  ]),
  mounted() {
    this.loadCommands();
  }
}
</script>
