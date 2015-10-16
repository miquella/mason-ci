import React from 'react'

import JobStore from '../stores/jobstore'
import JobList from './jobs/joblist'

export default class Jobs extends React.Component {
  constructor(props) {
    super(props);
    this.state = {jobs: JobStore.jobs};
    this.updateFromStore = this.updateFromStore.bind(this);
  }

  componentDidMount() {
    JobStore.addListener('change', this.updateFromStore);
  }

  componentWillUmount() {
    JobStore.removeListener('change', this.updateFromStore);
  }

  updateFromStore() {
    this.setState({jobs: JobStore.jobs});
  }

  render() {
    return(
      <div>
        <h1>Jobs</h1>
        <p></p>
        <JobList jobs={this.state.jobs} />
      </div>
    )
  }
}