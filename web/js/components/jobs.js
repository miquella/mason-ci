import React from 'react'

import JobStore from '../stores/JobStore'
import JobList from './jobs/JobList'

export default class Jobs extends React.Component {
  constructor(props) {
    super(props);
    this.state = {jobs: JobStore.jobs};
    this.updateFromStore = this.updateFromStore.bind(this);
  }

  componentDidMount() {
    JobStore.addListener('change_jobs', this.updateFromStore);
  }

  componentWillUmount() {
    JobStore.removeListener('change_jobs', this.updateFromStore);
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