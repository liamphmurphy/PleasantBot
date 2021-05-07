<template>
  <div class="home">
    
    <table>
      <tr>
        <th v-for="header in headers" :key="header"> {{ header }}</th>
      </tr>
      <template v-for="(point) in data">
        <tr
        v-for="(item, key) in point"
        :key="key">
        <td> {{key}} </td>
        <td v-for="attr in item" :key="attr"> {{attr}} </td>
        </tr>
      </template>
    </table>
  </div>
</template>

<script>

export default {
  name: 'Table',
  props: ['data', 'headers'],
}
</script>
