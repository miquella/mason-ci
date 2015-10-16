import React from 'react'
import JobListItem from './JobListItem'

export default (props) => (
  <table className="table table-striped">
    <thead>
      <tr>
        <th className="col-lg-3">Name</th>
        <th className="col-lg-8">Description</th>
        <th className="col-lg-1"></th>
      </tr>
    </thead>
    <tbody>
      {props.jobs.map((job) => {return(<JobListItem key={job.key} job={job} />)})}
    </tbody>
  </table>
)