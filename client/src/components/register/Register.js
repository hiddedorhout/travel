import React, {Component} from 'react'
import "./register.css"

class Register extends Component {
	render(){
		return(
			<div className="registerPage">
				<form className="registerForm">
					<input className="userName" type="text" placeholder="username" />
					<input type="email" className="email" placeholder="@email" />
					<input type="password" className="password" id="password1" placeholder="password" />
					<input type="password" className="password" id="password2" placeholder="repeat password" />
					<input type="submit" className="register" value="Register" />
					<span className="backlogin" onClick={e => this.props.setPage(e, "login")}>Back to login...</span>
				</form>
			</div>
		)
	}
}

export default Register