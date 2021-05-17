import UpdateTracker from './updatetracker';
import Moment from 'react-moment';  
import React from 'react';
import { Button } from 'react-bootstrap';
import * as moment from 'moment';
import requests from '../Requests.js';
  
class TrackerRow extends React.Component {
  constructor(props) {
    super(props);
  }

  renderEndDate(timestamp) {
    if (timestamp == "") {
      return (
        "Ongoing"
      )
    } else {
      return (
        <Moment format="YYYY/MM/DD hh:mm:ss">{timestamp}</Moment>
      )
    }
  }

  render() {

    const tracker = this.props.tracker;
    const timeStart = moment(tracker.start);
    let elapsed = "";
    let timeEnd = "";
    if (tracker.end === null) {
      elapsed = moment.duration(moment().utc().diff(timeStart));
    } else {
      timeEnd = moment(tracker.end);
      elapsed = moment.duration(timeEnd.diff(timeStart));
    }


    return(
      <tr>
        <td>{tracker.name}</td>
        <td><Moment format="YYYY/MM/DD hh:mm:ss">{tracker.start}</Moment></td>
        <td>{this.renderEndDate(timeEnd)}</td>
        <td>{elapsed.humanize()}</td>
        <td>
          <StopTracking tracker={tracker} refresh={this.props.refresh}></StopTracking>
        </td>
        <td>
          <DeleteTracker id={tracker.id} refresh={this.props.refresh}></DeleteTracker>
        </td>
        <td>
          <UpdateTracker tracker={tracker} refresh={this.props.refresh}></UpdateTracker>
        </td>
      </tr>
    )
  }
}


class StopTracking extends React.Component {
  constructor(props) {
    super(props);

    this.stopTracking = this.stopTracking.bind(this);

  }

  async stopTracking() {

    await requests.updateTracker({
      name: this.props.tracker.name,
      end: moment(),
      version: this.props.tracker.version
    }, this.props.tracker.id)

    await this.props.refresh();
  }

  render() {
    return (
      <Button variant="primary" disabled={this.props.tracker.end !== null} onClick={this.stopTracking.bind(this)}>Stop Tracking</Button>
    )
  }
}

class DeleteTracker extends React.Component {
  constructor(props) {
    super(props);

    this.deleteTracker = this.deleteTracker.bind(this);
  }

  async deleteTracker() {
    await requests.deleteTracker(this.props.id)

    await this.props.refresh();
  }

  render() {
    return (
      <Button variant="primary" onClick={this.deleteTracker.bind(this)}>Delete Tracker</Button>
    )
  }
}

export default TrackerRow;