import React from 'react';
import { Button, ButtonGroup } from 'react-bootstrap'

class Filter extends React.Component {

    constructor(props) {
        super(props);


    }

    render() {
        return (
            <ButtonGroup aria-label="Filters">
                <Button variant="secondary" onClick={ () => this.props.handler('all') } >All</Button>
                <Button variant="secondary" onClick={ () => this.props.handler('day') } >Day</Button>
                <Button variant="secondary" onClick={ () => this.props.handler('week') } >Week</Button>
                <Button variant="secondary" onClick={ () => this.props.handler('month') } >Month</Button>
            </ButtonGroup>
        )
    }
}

export default Filter;