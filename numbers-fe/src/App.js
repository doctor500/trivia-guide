import React from 'react';
import logo from './logo.svg';
import './App.css';
import axios from 'axios';

export default class PersonList extends React.Component {
  state = {
    number: '',
  }

  handleChange = event => {
    this.setState({ number: event.target.value });
  }

  handleSubmit = event => {
    event.preventDefault();

    var url = `https://faldyfin-numbers-be.herokuapp.com/numbersapi/${this.state.number}`
    axios(url, {
  auth: {
    username: "dora",
    password: "emon"
  },  withCredentials: true, method: 'get', })
      .then(res => {
        console.log(res);
        console.log(res.data);
        this.setState({trivia: res.data})
      })
  }

render() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Input the number to get their Trivia
        </p>
        <form onSubmit={this.handleSubmit}>
          <label>

            <input type="number" name="number" onChange={this.handleChange} />
          </label>
          <button type="submit">GO</button>
        </form>
        <br></br>
        {this.state.trivia}
        <a
          className="App-link"
          href="https://www.linkedin.com/in/faldy-findraddy/"
          target="_blank"
          rel="noopener noreferrer"
        >
          <br></br>
          <span role="img" aria-label="sheep">🐑</span>
        </a>
      </header>
    </div>
  );
}
}

//export default App;
