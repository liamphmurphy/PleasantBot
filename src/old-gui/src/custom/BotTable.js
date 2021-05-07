import React from 'react';
import BootstrapTable from 'react-bootstrap-table-next';
import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit';

class BotTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: true,
    }

    //this.getCommands = this.getCommands.bind(this);
    this.handleOnSelect = this.handleOnSelect.bind(this);
    this.handleOnSelectAll = this.handleOnSelectAll.bind(this);

    this.selectedRows = [] // stores any rows selected
  }


  //componentDidMount() {
  //  this.getCommands()
  //}


  // facilitates operations for selecting / deselecting rows and updating this.selectedRows
  handleOnSelect(row, isSelect) {
    if (isSelect){ // used to see if row is actually selected
      this.selectedRows.push(row.name) // add command to 
    } else {
      var index = this.selectedRows.indexOf(row.name)
      this.selectedRows.splice(index, 1)
    }

    this.props.rowData(this.selectedRows)
  }

  handleOnSelectAll(isSelect, rows) {
    if (isSelect) {
      // eslint-disable-next-line
      for (const [index, com] of rows.entries()) {
        console.log(com)
        this.selectedRows.push(com.name)
      }
      this.props.rowData(this.selectedRows)
    }
  }



  render() {
    const selectRow = {
      mode: 'checkbox',
      clickToSelect: true,
      selected: this.state.selected,
      onSelect: this.handleOnSelect,
      onSelectAll: this.handleOnSelectAll,
      style: { backgroundColor: "#007bff", color: "#FFFFFF" }
    };

    return (
        <ToolkitProvider bootstrap4 keyField="name" data={ this.props.data } columns={ this.props.columns } search>
          { props => (
            <Container fluid>
              <Row>
                <Col>
                  <b>{this.props.tableName}</b>
                  <hr />
                  <br />
                  <BootstrapTable { ...props.baseProps } boostrap4 selectRow={selectRow} pagination={ paginationFactory() }/>
                </Col>
              </Row>
            </Container>
        )}
        </ToolkitProvider>
    );
  
  }
}

export default BotTable;
