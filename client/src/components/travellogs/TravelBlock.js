import React, {Component} from 'react';
import './travelblock.css'
import moment from "moment";


class TravelBlock extends Component {

	state = {active: false, startTime: null, endTime: null, duration: null}

	setActive = () => {
		if (!this.state.active) {
			this.setState(prevState => ({ startTime: moment().format("DD-MM-YYYY HH:mm:ss"), active: !prevState.active}))
		}
		if (this.state.active && this.state.startTime !== null) {
			const end = moment().format("DD-MM-YYYY HH:mm:ss")
			const duration = moment.duration(moment(end, "DD-MM-YYYY HH:mm:ss") - moment(this.state.startTime, "DD-MM-YYYY HH:mm:ss"));
			this.setState({endTime: end, duration})
		}

		if (this.state.duration !== null) {
			this.setState(prevState => ({ active: !prevState.active, startTime: null, endTime: null, duration: null}))
		}

	}

	render(){
		return(
			<div className="block" id={this.state.active.toString()} onClick={this.setActive}>
				<div className="blockHeader">
					<div className="image" />
					{this.props.transport}
				</div>
				<div className="timeBlock">
					<div className="time">
						Start: {this.state.startTime}
					</div>
					<div className="time">
						End: {this.state.endTime}
					</div>
				</div>
			</div>
		)
	}
}

export default TravelBlock;