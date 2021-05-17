import React from 'react';
import { Form, Button } from 'react-bootstrap';
import requests from '../Requests.js';

class UpdateTracker extends React.Component {
    constructor(props) {
      super(props);
  
      this.state = {
        showForm: false,
        trackerName: ""
      }
  
      this.input = React.createRef();
  
      this.updateTracker = this.updateTracker.bind(this);
      this.handleChange = this.handleChange.bind(this);
      this.switchVisible = this.switchVisible.bind(this);
    }
  
    handleChange(event) {
      this.setState({
        trackerName: this.input.current.value
      })
    }
  
    switchVisible() {
      this.setState({
        showForm: !this.state.showForm
      })
    }
  
    async updateTracker() {
      await requests.updateTracker({
        name: this.state.trackerName,
        end: this.props.tracker.end,
        version: this.props.tracker.version
      }, this.props.tracker.id)
  
      await this.props.refresh();
    }
  
    renderEditForm() {
      if (this.state.showForm == true) {
        return (
          <Form>
            <Form.Group controlId="formTrackerName">
              <Form.Control type="text" placeholder="Enter new tracker name" ref={this.input} onChange={this.handleChange} />
              <Button variant="primary" onClick={this.updateTracker} >Update Tracker</Button>
            </Form.Group>
          </Form> 
        )
      }
    }
  
    render() {
      return (
        <div>
          <Button variant="primary" onClick={this.switchVisible}>Edit Tracker</Button>
          {this.renderEditForm()}
        </div>
      )
    }
  
}

export default UpdateTracker;