import React from 'react'
import JobStore from '../../stores/jobstore'

export default (props) => (
  <tr>
    <td>{props.job.name}</td>
    <td>{props.job.description}</td>
    <td><button className="btn btn-primary btn-xs" onClick={JobStore.createBuild.bind(JobStore, props.job.key)}>Run Build</button></td>
  </tr>
)