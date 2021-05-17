import React from 'react';
import { Form, Button } from 'react-bootstrap';
import requests from '../Requests.js';
import * as moment from 'moment';

class NewTracker extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            trackerName: ""
        };

        this.input = React.createRef();

        this.handleChange = this.handleChange.bind(this);
        this.createNewTracker = this.createNewTracker.bind(this);
    }

    async createNewTracker() {

        await requests.createTracker({
            start:  moment(),
            name:   this.state.trackerName
        });

        await this.props.refresh();
    }

    handleChange(event) {
        this.setState({
            trackerName: this.input.current.value
        })
    }

    render() {
        return (
            <Form>
                <Form.Group controlId="formTrackerName">
                    <Form.Label>Tracker Name: </Form.Label>
                    <Form.Control type="text" placeholder="Enter new tracker name" ref={this.input} onChange={this.handleChange} />
                    <Button variant="primary" onClick={this.createNewTracker} >Create Tracker</Button>
                </Form.Group>
            </Form>
        );
    }
}

export default NewTracker;