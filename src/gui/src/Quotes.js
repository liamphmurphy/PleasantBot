import React from 'react';
import BootstrapTable from 'react-bootstrap-table-next';
import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit';

const { SearchBar } = Search;

const axios = require('axios').default;

class Quotes extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      quotes: []
    }

    this.getQuotes = this.getQuotes.bind(this);
    this.handleOnSelect = this.handleOnSelect.bind(this);
    this.handleOnSelectAll = this.handleOnSelectAll.bind(this);

    this.selectedRows = [] // stores any rows selected
  }


  componentDidMount() {
    this.getQuotes()
  }

  // used to reload the commands table when operations are performed
  getQuotes() {
    axios.get("http://" + window.location.hostname + ":8080/getquotes").then((response) => {
      this.setState({quotes: response.data})
    }).catch(function(error) {
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
    const {quotes} = this.state;

    // turns object into an array, easier to manipulate
    console.log(quotes)
    const arr = []

    // turns object into an array, easier to manipulate
    Object.keys(quotes).forEach(key => arr.push({id: key, value: quotes[key]}))
    const columns = [{dataField: 'id', text: 'ID', sort: true}, 
                    {dataField: 'value.Quote', text: 'Quote'},
                    {dataField: 'value.Timestamp', text: 'Timestamp'},
                    {dataField: 'value.Submitter', text: 'Submitter'} ];

    const selectRow = {
      mode: 'checkbox',
      clickToSelect: true,
      selected: this.state.selected,
      onSelect: this.handleOnSelect,
      onSelectAll: this.handleOnSelectAll,
      style: { backgroundColor: "#007bff", color: "#FFFFFF" }
    };

    

    return (
        <ToolkitProvider bootstrap4 keyField="id" data={ arr } columns={ columns } search>
          { props => (
            <Container fluid>
              <Row>
                <Col>
                  <b>Quotes</b>
                  <hr />
                  <Row>
                    <Col md={8}>
                      <div className="search">
                        <SearchBar { ...props.searchProps } />
                      </div>
                    </Col>
                    <Col md={4}>
                      
                    </Col>
                  </Row>
                  <br />
                  <BootstrapTable { ...props.baseProps } boostrap4 selectRow={selectRow} pagination={ paginationFactory() }/>
                </Col>
              </Row>
              <Row>
                <Col>
                  
                </Col>
              </Row>
            </Container>
        )}
        </ToolkitProvider>
    );
  }
}

export default Quotes;
