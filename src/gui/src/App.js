import React from 'react';

// application page imports
import Commands from './Commands';
import Help from './Help';

// bootstrap imports
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import 'bootstrap/dist/css/bootstrap.min.css';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Container from 'react-bootstrap/Container';

function App() {
  return (
    <div>
      <Router>
        <Navbar collapseOnSelect bg="primary" variant="dark" fixed="top" expand="lg"> 
          <Navbar.Brand href="/home">PleasantBot</Navbar.Brand>
          <Navbar.Toggle aria-controls="responsive-navbar-nav" />
          <Navbar.Collapse id="responsive-navbar-nav">
            <Nav className="mr-auto">
              <Nav.Link href="/home">Home</Nav.Link>
              <Nav.Link href="/commands">Commands</Nav.Link>
            </Nav>

            <Nav>
              <Nav.Link href="/help">Help</Nav.Link>
            </Nav>
          </Navbar.Collapse>
        </Navbar>
        <br /><br /><br />
        <Container>
          <Switch>
            <Route exact path="/commands">
              <Commands />
            </Route>

            <Route exact path="/help">
              <Help />
            </Route>
          </Switch>
        </Container>
      </Router>
    </div>
  );
}

export default App;
