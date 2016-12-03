const URL = 'ws:localhost:5000/';

export default class API {
	constructor () {
		this.sock = new WebSocket(URL)
		this.sock.onmessage = this.onMessage;
	}

	onMessage() {
		
	}
}
