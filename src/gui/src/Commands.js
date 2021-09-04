import React from 'react';
import BootstrapTable from 'react-bootstrap-table-next';
import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import AddCommand from './custom/AddCommand';
import DeleteCommand from './custom/DeleteCommand';

import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit';

const { SearchBar } = Search;

const axios = require('axios').default;

class Commands extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: true,
      commands: []
    }

    this.getCommands = this.getCommands.bind(this);
    this.handleOnSelect = this.handleOnSelect.bind(this);
    this.handleOnSelectAll = this.handleOnSelectAll.bind(this);

    this.selectedRows = [] // stores any rows selected
  }


  componentDidMount() {
    this.getCommands()
  }

  // used to reload the commands table when operations are performed
  getCommands() {
    axios.get("http://" + window.location.hostname + ":8080/getcoms").then((response) => {
      this.setState({loaded: true, commands: response.data})
    }).catch(function(error) {
      this.setState({loaded: false})
      console.log(error)
    });
  }

  // facilitates operations for selecting / deselecting rows and updating this.selectedRows
  handleOnSelect(row, isSelect) {
    if (isSelect){ // used to see if row is actually selected
      this.selectedRows.push(row.name) // add command to 
    } else {
      var index = this.selectedRows.indexOf(row.name)
      this.selectedRows.splice(index, 1)
    }
  }

  handleOnSelectAll(isSelect, rows) {
    if (isSelect) {
      for (const [index, com] of rows.entries()) {
        console.log(com)
        this.selectedRows.push(com.name)
      }
    }
  }



  render() {
    const {loaded, commands} = this.state;

    const arr = []

    // turns object into an array, easier to manipulate
    Object.keys(commands).forEach(key => arr.push({name: key, value: commands[key]}))
    const columns = [{dataField: 'name', text: 'Command Name', sort: true}, 
                    {dataField: 'value.Response', text: 'Command Response'},
                    {dataField: 'value.Perm', text: 'Command Permissions'}];

    const selectRow = {
      mode: 'checkbox',
      clickToSelect: true,
      selected: this.state.selected,
      onSelect: this.handleOnSelect,
      onSelectAll: this.handleOnSelectAll,
      style: { backgroundColor: "#007bff", color: "#FFFFFF" }
    };

    

    if (loaded) { // indicates that the bot may not be running
      return (
          <ToolkitProvider bootstrap4 keyField="name" data={ arr } columns={ columns } search>
            { props => (
              <Container fluid>
                <Row>
                  <Col>
                    <b>Commands</b>
                    <hr />
                    <Row>
                      <Col md={8}>
                        <div className="search">
                          <SearchBar { ...props.searchProps } />
                        </div>
                      </Col>
                      <Col md={4}>
                        <DeleteCommand selectedRows={ this.selectedRows } getCommands={ this.getCommands.bind(this) }/>
                      </Col>
                    </Row>
                    <br />
                    <BootstrapTable { ...props.baseProps } boostrap4 selectRow={selectRow} pagination={ paginationFactory() }/>
                  </Col>
                </Row>
                <Row>
                  <Col>
                    <AddCommand getCommands={ this.getCommands.bind(this) }/>
                  </Col>
                </Row>
              </Container>
          )}
          </ToolkitProvider>
      );
    } else {
      return <h1>ERROR. Please ensure the bot is running.</h1>;
    }
  }
}

export default Commands;
