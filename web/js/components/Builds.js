import React from 'react'

import JobStore from '../stores/JobStore'
import BuildList from './builds/BuildList'

export default class Builds extends React.Component {
  constructor(props) {
    super(props);
    this.state = {job: props.params.job, builds: JobStore.builds};
    this.updateFromStore = this.updateFromStore.bind(this);
    JobStore.getBuilds(props.params.job);
  }

  componentDidMount() {
    JobStore.addListener('change_builds', this.updateFromStore);
  }

  componentWillUmount() {
    JobStore.removeListener('change_builds', this.updateFromStore);
  }

  updateFromStore() {
    this.setState({builds: JobStore.builds});
  }

  render() {
    return(
      <div>
        <h1>Builds for {this.state.job}</h1>
        <p></p>
        <BuildList builds={this.state.builds} />
      </div>
    )
  }
}