var http = require("http");


const postIt = () => {
	var options = {
		hostname: 'www.postcatcher.in',
		port: 80,
		path: '/catchers/544b09b4599c1d0200000289',
		method: 'POST',
		headers: {
				'Content-Type': 'application/json',
		}
	};
	var req = http.request(options, function(res) {
		console.log('Status: ' + res.statusCode);
		console.log('Headers: ' + JSON.stringify(res.headers));
		res.setEncoding('utf8');
		res.on('data', function (body) {
			console.log('Body: ' + body);
		});
	});
	req.on('error', function(e) {
		console.log('problem with request: ' + e.message);
	});
	// write data to request body
	req.write('{"string": "Hello, World"}');
	req.end();
}

var server = http.createServer(function(req, resp) {
  resp.writeHead(200, {"Content-Type": "application/json"});
  if (req.method === "POST") {

    let body = '';
    req.on('data', chunk => {
        body += chunk.toString(); // convert Buffer to string
    });
    req.on('end', () => {
			console.log(body);
			const obj = { msg: 'ola' };
			resp.write(JSON.stringify(obj));
			resp.end();
    });

    console.log('REQ', req.data);
  } else {
		const obj = { msg: 'ola' };
		resp.write(JSON.stringify(obj));
		resp.end();
  }

});


server.listen(8080);
console.log("Server is listening");
