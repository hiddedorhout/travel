import React, {Component} from 'react'
import "./login.css"

class Login extends Component {
	render() {
		return(
			<div className="loginPage">
				<form className="loginForm">
					<input type="text" className="userName" placeholder="username" />
					<input type="password" className="password" placeholder="password" />
					<input type="submit" className="login" value="Login"/>
					<span className="notRegistered" onClick={e => this.props.setPage(e, "register")}>Not registered?</span>
				</form>
			</div>
		)
	}
}

export default Login;