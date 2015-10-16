import React from 'react'
import { Link } from 'react-router'

import JobStore from '../../stores/JobStore'

export default (props) => (
  <tr>
    <td><Link to={'/jobs/' + props.job.key + '/builds'}>{props.job.name}</Link></td>
    <td>{props.job.description}</td>
    <td><button className="btn btn-primary btn-xs" onClick={JobStore.createBuild.bind(JobStore, props.job.key)}>Run Build</button></td>
  </tr>
)