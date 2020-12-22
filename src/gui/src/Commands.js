import React from 'react';
import BootstrapTable from 'react-bootstrap-table-next';
import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';

const axios = require('axios').default;

class Commands extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: true,
      commands: []
    }
  }


  componentDidMount() {
    axios.get("http://localhost:8080/getcoms").then((response) => {
      this.setState({loaded: true, commands: response.data})
    }).catch(function(error) {
      this.setState({loaded: false})
      console.log(error)
    });
  }

  render() {
    const {loaded, commands} = this.state;
    //let arr = [];

    const arr = []
    Object.keys(commands).forEach(key => arr.push({name: key, value: commands[key]}))
    const columns = [{dataField: 'name', text: 'Command Name', sort: true}, 
                    {dataField: 'value.Response', text: 'Command Response'},
                    {dataField: 'value.ModeratorPerms', text: 'Command Permissions'}];

    const selectRow = {
      mode: 'checkbox',
      clickToSelect: true,
      selected: this.state.selected,
      onSelect: this.handleOnSelect,
      onSelectAll: this.handleOnSelectAll
    };


    if (loaded) { // indicates that the bot may not be running
      return (
        <div>
          <BootstrapTable bootstrap4 keyField='name' data={arr} columns={columns} selectRow={selectRow} />
        </div>
      );
    } else {
      return <h1>ERROR. Please ensure the bot is running.</h1>;
    }
  }
}

export default Commands;
