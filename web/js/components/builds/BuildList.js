import React from 'react'
import BuildListItem from './BuildListItem'

export default (props) => (
  <table className="table table-striped">
    <thead>
      <tr>
        <th className="col-lg-3">Build Number</th>
        <th className="col-lg-8">Job Name</th>
        <th className="col-lg-1">Status</th>
      </tr>
    </thead>
    <tbody>
      {props.builds.map((build) => {return(<BuildListItem key={build.key} build={build} />)})}
    </tbody>
  </table>
)