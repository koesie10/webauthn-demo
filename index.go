package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return c.HTML(http.StatusOK, `<html>
<head>
	<title>WebAuthn demo</title>
	<style>
		.hide {
			display: none;
		}
	</style>

	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
</head>
<body>
<div class="container mt-4">
	<h1>WebAuthn Demo</h1>

	<p>This is a demo of the WebAuthn library for Go.</p>
	<p>You can try registering or logging in. If the username is not found when logging in, you can login as any account where your authenticator is registered. So you might be logged in to the account you registered before.</p>

	<div class="card mt-3">
		<div class="card-body">
			<p class="hide lead" id="registerLoading">Registering... Please tap your authenticator.</p>

			<form method="post" id="registerForm">
				<input type="text" name="name" id="registerName" class="form-control" placeholder="Username" />
				<button type="submit" class="btn btn-success mt-3">Register</button>
			</form>
		</div>
	</div>

	<div class="card mt-3">
		<div class="card-body">
			<p class="hide lead" id="loginLoading">Logging in... Please tap your authenticator.</p>

			<form method="post" id="loginForm">
				<input type="text" name="name" id="loginName" class="form-control" placeholder="Username" />
				<button type="submit" class="btn btn-primary mt-3">Login</button>
			</form>
		</div>
	</div>
</div>

	<script type="text/javascript">
// This is a modification of the example class, where the URLs have been changed to include the name.
class WebAuthn {
	// Decode a base64 string into a Uint8Array.
	static _decodeBuffer(value) {
		return Uint8Array.from(atob(value), c => c.charCodeAt(0));
	}

	// Encode an ArrayBuffer into a base64 string.
	static _encodeBuffer(value) {
		return btoa(new Uint8Array(value).reduce((s, byte) => s + String.fromCharCode(byte), ''));
	}

	// Checks whether the status returned matches the status given.
	static _checkStatus(status) {
		return res => {
			if (res.status === status) {
				return res;
			}
			throw new Error(res.statusText);
		};
	}

	register(name) {
		return fetch('/webauthn/registration/start/' + name, {
				method: 'POST'
			})
			.then(WebAuthn._checkStatus(200))
			.then(res => res.json())
			.then(res => {
				res.publicKey.challenge = WebAuthn._decodeBuffer(res.publicKey.challenge);
				res.publicKey.user.id = WebAuthn._decodeBuffer(res.publicKey.user.id);
				if (res.publicKey.excludeCredentials) {
					for (var i = 0; i < res.publicKey.excludeCredentials.length; i++) {
						res.publicKey.excludeCredentials[i].id = WebAuthn._decodeBuffer(res.publicKey.excludeCredentials[i].id);
					}
				}
				return res;
			})
			.then(res => navigator.credentials.create(res))
			.then(credential => {
				return fetch('/webauthn/registration/finish/' + name, {
					method: 'POST',
					headers: {
						'Accept': 'application/json',
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						id: credential.id,
						rawId: WebAuthn._encodeBuffer(credential.rawId),
						response: {
							attestationObject: WebAuthn._encodeBuffer(credential.response.attestationObject),
							clientDataJSON: WebAuthn._encodeBuffer(credential.response.clientDataJSON)
						},
						type: credential.type
					}),
				})
			})
			.then(WebAuthn._checkStatus(201));
	}

	login(name) {
		return fetch('/webauthn/login/start/' + name, {
				method: 'POST'
			})
			.then(WebAuthn._checkStatus(200))
			.then(res => res.json())
			.then(res => {
				res.publicKey.challenge = WebAuthn._decodeBuffer(res.publicKey.challenge);
				if (res.publicKey.allowCredentials) {
					for (let i = 0; i < res.publicKey.allowCredentials.length; i++) {
						res.publicKey.allowCredentials[i].id = WebAuthn._decodeBuffer(res.publicKey.allowCredentials[i].id);
					}
				}
				return res;
			})
			.then(res => navigator.credentials.get(res))
			.then(credential => {
				return fetch('/webauthn/login/finish/' + name, {
					method: 'POST',
					headers: {
						'Accept': 'application/json',
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						id: credential.id,
						rawId: WebAuthn._encodeBuffer(credential.rawId),
						response: {
							clientDataJSON: WebAuthn._encodeBuffer(credential.response.clientDataJSON),
							authenticatorData: WebAuthn._encodeBuffer(credential.response.authenticatorData),
							signature: WebAuthn._encodeBuffer(credential.response.signature),
							userHandle: WebAuthn._encodeBuffer(credential.response.userHandle),
						},
						type: credential.type
					}),
				})
			})
			.then(WebAuthn._checkStatus(200));
	}
}

let registerPending = false;
let loginPending = false;

let w = new WebAuthn();

document.getElementById("registerForm").onsubmit = function(e) {
	e.preventDefault();

	if (registerPending) return;
	registerPending = true;

	document.getElementById("registerLoading").classList.remove("hide");

	const name = document.getElementById("registerName").value;
	w.register(name)
			.then(res => alert('This authenticator has been registered'))
			.catch(err => {
				console.error(err)
				alert('Failed to register: ' + err);
			})
			.then(() => {
				registerPending = false;
				document.getElementById("registerLoading").classList.add("hide");
			});
};

document.getElementById("loginForm").onsubmit = function(e) {
	e.preventDefault();

	if (loginPending) return;
	loginPending = true;

	document.getElementById("loginLoading").classList.remove("hide");

	const name = document.getElementById("loginName").value;
	w.login(name)
			.then(res => res.json())
			.then(res => alert('You have been logged in to ' + res.name))
			.catch(err => {
				console.error(err)
				alert('Failed to login: ' + err);
			})
			.then(() => {
				loginPending = false;
				document.getElementById("loginLoading").classList.add("hide");
			});
};
	</script>
</body>
</html>`)
}
