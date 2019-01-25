import React, { Component } from 'react'
import './travellogs.css'
import Block from "./TravelBlock.js"
import moment from 'moment';

class TravelLogs extends Component {

	state = {sessionID:"" , active: false}

	render(){
		return(
			<div className="travelLogPage">
				<div className="date">
					{moment().format("dddd Do of MMMM YYYY",)}
				</div>
				<ul className="modeOfTransportList">
				{
					["Train", "Bus", "Bike"].map( (modeOfTransport, i) => {
						return(
							<li key={i} className="transportType">
								<Block transport={modeOfTransport} index={i}/>
							</li>		
						)
					})
				}
				</ul>
			</div>
		)
	}
}

export default TravelLogs;