import React from 'react';
import BootstrapTable from 'react-bootstrap-table-next';

import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';
import 'react-bootstrap-table2-paginator/dist/react-bootstrap-table2-paginator.min.css';
import paginationFactory from 'react-bootstrap-table2-paginator';

import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

const axios = require('axios').default;

class Dashboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loaded: true,
      bans: [],
      stats: []
    }
  }


  componentDidMount() {
    axios.get("http://" + window.location.hostname + ":8080/getbanhistory").then((response) => {
      this.setState({loaded: true, bans: response.data})
    }).catch(function(error) {
      //this.setState({loaded: false})
      console.log(error)
    });

    axios.get("http://" + window.location.hostname + ":8080/getstats").then((response) => {
      this.setState({stats: response.data})
    }).catch(function(error) {
      //this.setState({loaded: false})
      console.log(error)
    });
  }

  render() {
    const {loaded, bans} = this.state;

    const columns = [{dataField: 'User', text: 'Username'}, 
                    {dataField: 'Reason', text: 'Reason'},
                    {dataField: 'Timestamp', text: 'Timestamp'}];
    if (loaded) { // indicates that the bot may not be running
      return (
        <Container fluid>
          <Row>
            <Col>
              <b>Quick Stats</b>
              <p>Commands: {this.state.stats.Commands}</p>
              <p>Quotes: {this.state.stats.Quotes}</p>
              <p>Bans: {this.state.stats.Bans}</p>
            </Col>
            <Col>
                <b>Top Command</b>
                <br />
                <p>{this.state.stats.TopCommand}: {this.state.stats.TopComCount}</p>

                <b>Top Chatter</b>
                <br />
                <p>{this.state.stats.TopChatter}: {this.state.stats.TopChatCount}</p>
            </Col>
          </Row>
          <Row>
          <Col>
              <b>Ban History</b>
              <BootstrapTable bootstrap4 keyField='name' data={bans} columns={columns} pagination={ paginationFactory() }/>
            </Col>
          </Row>
        </Container>
      );
    } else {
      return <h1>ERROR! Is the bot running?</h1>
    }
  }
}

export default Dashboard;
