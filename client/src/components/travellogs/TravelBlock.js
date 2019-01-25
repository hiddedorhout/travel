import React, {Component} from 'react';
import './travelblock.css'
import moment from "moment";
import { timeout } from 'q';


class TravelBlock extends Component {

	state = {active: false, startTime: null, endTime: null, duration: null}

	setActive = () => {
		if (!this.state.active) {
			this.setState(prevState => ({ startTime: Date.now(), active: !prevState.active}))
		}
		if (this.state.active && this.state.startTime !== null) {
			const endTime = Date.now()
			const duration = endTime - this.state.startTime;
			this.setState({endTime, duration})
		}

		if (this.state.duration !== null) {
			this.setState(prevState => ({ active: !prevState.active, startTime: null, endTime: null, duration: null}))
		}
	}

	addMinutes = (e, stateName) => {
		e.preventDefault();
		const newTime = this.state[stateName] + (60*1000)
		this.setState({[stateName]: newTime});
		if (stateName === "startTime") {
			const newDuration = this.state.endTime - newTime
			if (newDuration > 0) {
				this.setState({duration: newDuration})
			}
		}
		if (stateName === "endTime"){
			const newDuration = newTime - this.state.startTime
			if (newDuration > 0) {
				this.setState({duration: newDuration})
			}
		}
	}

	minMinutes = (e, stateName) => {
		e.preventDefault();
		const newTime = this.state[stateName] - (60 * 1000)
		this.setState({[stateName]: newTime});
		if (stateName === "startTime") {
			const newDuration = this.state.endTime - newTime
			if (newDuration > 0) {
				this.setState({duration: newDuration})
			}
		}
		if (stateName === "endTime"){
			const newDuration = newTime - this.state.startTime
			if (newDuration > 0) {
				this.setState({duration: newDuration})
			}
		}
	}

	render(){
		return(
			<div className="main">
				<div className="block">
					<div className="blockHeader">
					<div className="image" />
					{this.props.transport}
					</div>
					<div className="timeBlock">
					{
						this.state.startTime !== null &&
						<div className="timeDisplay">
							<button className="manipulateTimeButton" onClick={e => this.minMinutes(e, "startTime")}>-1</button>
							<div className="time">
							Start: {moment(this.state.startTime).format("HH:mm")}
							</div>
							<button className="manipulateTimeButton" onClick={e => this.addMinutes(e, "startTime")}>+1</button>
						</div>
					}
					{
						this.state.endTime !== null &&
						<div className="timeDisplay">
							<button className="manipulateTimeButton" onClick={e => this.minMinutes(e, "endTime")}>-1</button>
							<div className="time">
							End: {moment(this.state.endTime).format("HH:mm")}
							</div>
							<button className="manipulateTimeButton" onClick={e => this.addMinutes(e, "endTime")}>+1</button>
						</div>
					}
					{
						this.state.duration !== null &&
						<div className="time">
						Duration: {moment.duration(this.state.duration).minutes()} minutes
						</div>
					}
					</div>
				</div>
				<div className="starterStopper" id={this.state.active.toString()} onClick={this.setActive}>
					{
						!this.state.active && 
							<span className="time">Start!</span>
					}
					{
						this.state.active && this.state.endTime === null &&
						<span className="time">Running</span>
					}
					{
						this.state.endTime !== null &&
							<span className="time">Stopped</span>
					}
				</div>
			</div>
		)
	}
}

export default TravelBlock;