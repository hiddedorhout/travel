import React, { Component } from 'react';
import Login from "./components/login/Login"
import Register from "./components/register/Register"
import TravelLogs from "./components/travellogs/Travellogs"
import "./app.css"

class App extends Component {

  state = {nav: "travel"}

  setPage = (e, page) => {
    e.preventDefault();
    this.setState({nav: page});
  }

  render() {
    return (
      <div className="App">
        <div className="header">
          Travel
        </div>
        <div className="pages">
          {
            this.state.nav === "login" &&
            <Login setPage={this.setPage}/>
          }
          {
            this.state.nav === "register" && 
            <Register setPage={this.setPage}/>
          }
          {
            this.state.nav === "travel" &&
            <TravelLogs setPage={this.setPage}/>
          }
        </div>
      </div>
    );
  }
}

export default App;
