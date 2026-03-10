const http = require('http');

const options = {
	hostname: '127.0.0.1',
	port: 8080,
	path: '/materiais/favoritos?usuario_id=2',
	method: 'GET'
};

const req = http.request(options, (res) => {
	let data = '';
	res.on('data', (chunk) => {
		data += chunk;
	});
	res.on('end', () => {
		console.log("Status Code:", res.statusCode);
		fmtJson(data);
	});
});

req.on('error', (e) => {
	console.error(`problem with request: ${e.message}`);
});

req.end();

function fmtJson(data) {
	try {
		const obj = JSON.parse(data);
		console.log("Response Body:", JSON.stringify(obj, null, 2));
	} catch (e) {
		console.log("Raw Response:", data);
	}
}
