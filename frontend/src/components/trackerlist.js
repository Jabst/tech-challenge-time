import React from 'react';
import { Table } from 'react-bootstrap';
import TrackerRow from './trackerrow';

class TrackerList extends React.Component {

    constructor(props) {
      super(props);
  
    }
  
    render() {
  
      let trackerRows;
  
      if (this.props.trackers.length === 0) {
  
      } else {
        trackerRows = this.props.trackers.map( elem => (
          <TrackerRow key={elem.id} tracker={elem} refresh={this.props.refresh} />
        ));
      }
  
      return (
      <Table striped bordered hover size="sm">
        <thead>
          <tr>
            <th>Name</th>
            <th>Started At</th>
            <th>Ended At</th>
            <th>Time Elapsed</th>
            <th></th>
            <th></th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {trackerRows}
        </tbody>
      </Table>
      );
    }
}

export default TrackerList;