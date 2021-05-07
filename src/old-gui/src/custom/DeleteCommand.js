import React from 'react';
import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';

import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';


const axios = require('axios').default;

class DeleteCommand extends React.Component {
  constructor(props) {
    super(props);

    this.handleDelete = this.handleDelete.bind(this)
  }

  // send all selected row command names to back-end
  handleDelete(e) {
    e.preventDefault()

    let config = {
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
      }
    }

    axios.post("http://" + window.location.hostname + ":8080/delcom", JSON.stringify(this.props.selectedRows), config).then((res) => {
        console.log(res.data)
      }).catch((error) => {
          console.log(error)
      });

    
    this.props.getCommands() // parent function to update commands table
  }

  render() {
    return (
        <Container>
            <Form onSubmit={this.handleDelete}>
                <Button variant="danger" as="input" type="submit" value="Delete Command" />{' '}
            </Form>
        </Container>
    );
    
  }
}

export default DeleteCommand;