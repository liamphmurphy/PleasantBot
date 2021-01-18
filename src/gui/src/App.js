import React from 'react';

// application page imports
import Commands from './Commands';
import Help from './Help';
import Dashboard from './Dashboard';
import Quotes from './Quotes';
import Login from './Login';

// bootstrap imports
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';

// services
import {authenticated, sendToken} from './services/Auth'

import 'bootstrap/dist/css/bootstrap.min.css';
import './css/style.css';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect
} from "react-router-dom";
import Container from 'react-bootstrap/Container';


// renders the nav bar
function renderNavBar() {
  return (
    <Navbar collapseOnSelect bg="dark" variant="dark" fixed="top" expand="lg"> 
          <Navbar.Brand href="/">PleasantBot</Navbar.Brand>
          <Navbar.Toggle aria-controls="responsive-navbar-nav" />
          <Navbar.Collapse id="responsive-navbar-nav">
            <Nav className="mr-auto">
              <Nav.Link href="/">Home</Nav.Link>
              <Nav.Link href="/commands">Commands</Nav.Link>
              <Nav.Link href="/quotes">Quotes</Nav.Link>
              <Nav.Link href="/timers">Timers</Nav.Link>
            </Nav>

            <Nav>
              <Nav.Link href="/login">Login</Nav.Link>
              <Nav.Link href="/help">Help</Nav.Link>
            </Nav>
          </Navbar.Collapse>
        </Navbar>
  )
}

function App() {
  if (document.location.hash != "") { 
    sendToken(document.location.hash.substr(1)); // send new oauth token to backend
  }

  return (
    <div>
      <Router>
        {renderNavBar()}
        <br /><br /><br />
        <Container>
          <Switch>
            <Route exact path="/login">
              <Login />
            </Route>
            <Route exact path="/">
              {authenticated() ? <Dashboard /> : <Redirect to="/login" />}
            </Route>
            <Route exact path="/commands">
              <Commands />
            </Route>
            <Route exact path="/quotes">
              <Quotes />
            </Route>
            <Route exact path="/timers">

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
