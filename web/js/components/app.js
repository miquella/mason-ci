import React from 'react'
import { Link, IndexLink } from 'react-router'

export default class App extends React.Component {
  render() {
    return (
      <div style={{paddingTop: 70 + 'px'}}>
        <nav className="navbar navbar-default navbar-fixed-top">
          <div className="container">
            <div className="navbar-header">
              <span className="navbar-brand">Mason-CI</span>
            </div>
            <div className="navbar-collapse collapse">
              <ul className="nav navbar-nav">
                <li><IndexLink to="/">Home</IndexLink></li>
                <li><Link to="/jobs">Jobs</Link></li>
              </ul>
            </div>
          </div>
        </nav>
        <div className="container">
          {this.props.children}
        </div>
      </div>
    )
  }
}