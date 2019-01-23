import React, { Component } from 'react'
import './travellogs.css'
import Block from "./TravelBlock.js"

class TravelLogs extends Component {
	render(){
		return(
			<div className="travelLogPage">
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