import React, {Component} from 'react';
import './travelblock.css'


class TravelBlock extends Component {

	state = {active: false, startTime: null, endTime: null}

	setActive = () => {
		if (!this.state.active) {
			this.setState({startTime: Date.now()})
		}
		if (this.state.active && this.state.startTime !== null) {
			this.setState({endTime: Date.now()})
		}

		this.setState(prevState => ({active: !prevState.active}))
	}

	render(){
		return(
			<div className="block" onClick={this.setActive}>
				<div className="blockHeader">
					<div className="image" />
					{this.props.transport}
				</div>
				{this.state.startTime}
				{this.state.endTime}
			</div>
		)
	}
}

export default TravelBlock;