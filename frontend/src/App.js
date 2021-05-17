import './App.css';
import requests from './Requests.js';
import React from 'react';
import Container from 'react-bootstrap/Container';
import * as moment from 'moment';
import NewTracker from './components/newtracker';
import Filter from './components/filter';
import TrackerList from './components/trackerlist'
import { render } from '@testing-library/react';


class App extends React.Component {
  constructor(props) {
    super(props);

    this.state = { 
      trackers: [],
      timeFilter: 'all'
    };

    
    this.refresh = this.refresh.bind(this);
    this.filterHandler = this.filterHandler.bind(this);
    this.trackersLabel = this.trackersLabel.bind(this);
    
  }
  
  async componentDidMount() {
    await this.refresh();
  }

  async refresh() {
    let trackers = [];

    switch (this.state.timeFilter) {
      case "all":
        trackers = await requests.getTrackers()
        break;
      case "day":
        trackers = await requests.getTrackers(moment().startOf('day').toISOString(), moment().endOf('day').toISOString())
        break;
      case "week":
        trackers = await requests.getTrackers(moment().startOf('week').isoWeekday(1).toISOString(),moment().endOf('week').isoWeekday(1).toISOString())
        break;
      case "month": 
        trackers = await requests.getTrackers(moment().startOf('month').toISOString(),moment().endOf('month').isoWeekday(1).toISOString())
        break;
    }

    this.setState({
      trackers: trackers
    });
  }

  filterHandler(arg) {
    this.setState({
      timeFilter: arg
    }, () => {
      this.refresh();
    });
  }

  trackersLabel() {
    switch (this.state.timeFilter) {
      case 'all':
        return (
          <h2>All trackers</h2>
        )
      default:
        return(
          <h2>Trackers listed by {this.state.timeFilter}</h2>
        )
    }
  }
  
  render() {
    return (
      <Container>
        <h1>Trackers</h1>
        <h2>Filter by</h2>
        <Filter handler={this.filterHandler}></Filter>
        <h2>Create tracker</h2>
        <NewTracker refresh={this.refresh} ></NewTracker>
        {this.trackersLabel()}
        <TrackerList refresh={this.refresh} trackers={this.state.trackers}></TrackerList>
      </Container>
        
    );
  }

}

export default App;
