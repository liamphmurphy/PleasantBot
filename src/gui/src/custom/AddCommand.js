import React from 'react';
import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';

import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import Col from 'react-bootstrap/Col';




const axios = require('axios').default;

class AddCommand extends React.Component {
  constructor(props) {
    super(props);

    this.onChangeComName = this.onChangeComName.bind(this);
    this.onChangeComResponse = this.onChangeComResponse.bind(this);
    this.onChangeComPermission = this.onChangeComPermission.bind(this);
    this.onSubmit = this.onSubmit.bind(this);

    this.state = {
      name: "",
      response: "",
      perm: "all"
    }

  }
  


  // update com name from form
  onChangeComName(e) {
    this.setState({name: e.target.value})
  }

  // update com response from form
  onChangeComResponse(e) {
    this.setState({response: e.target.value})
  }

  // update com permission from form
  onChangeComPermission(e) {
    this.setState({perm: e.target.value})
  }

  onSubmit(e) {
      e.preventDefault()

    let config = {
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
      }
    }
    
      // prepare command object
      const command = {
          CommandName: this.state.name,
          Response: this.state.response,
          Perm: this.state.perm
      }


      axios.post("http://" + window.location.hostname + ":8080/addcom", JSON.stringify(command), config).then((res) => {
        console.log(res.data)
      }).catch((error) => {
          console.log(error)
      });


      // return to default state
      this.setState({
        name: "",
        response: "",
        perm: "all"
      })

      this.props.getCommands() // parent function to update commands table
  }

  

  render() {
    return (
        <Container>
            <b>Add Command</b>
            <Form onSubmit={this.onSubmit}>
              <Form.Row>
                  <Form.Group as={Col} md="4" controlId = "formName">
                    <Form.Label>Command Name</Form.Label>
                    <Form.Control required type="text" value={this.state.name} onChange={this.onChangeComName} />
                  </Form.Group>

                  <Form.Group as={Col} md="4" controlId = "formResponse">
                    <Form.Label>Command Response</Form.Label>
                    <Form.Control required type="text" value={this.state.response} onChange={this.onChangeComResponse} />
                  </Form.Group>

                  <Form.Group as={Col} md="4" controlId = "formResponse">
                    <Form.Label>Command Permission</Form.Label>
                    <Form.Control required custom as="select" value={this.state.perm} onChange={this.onChangeComPermission}>
                      <option value="all">All</option>
                      <option value="subscriber">Subscriber</option>
                      <option value="moderator">Moderator</option>
                      <option value="broadcaster">Broadcaster</option>
                    </Form.Control>
                  </Form.Group>
                </Form.Row>
              <Button as="input" type="submit" value="Create Command" />{' '}
            </Form>
        </Container>
    );
    
  }
}

export default AddCommand;