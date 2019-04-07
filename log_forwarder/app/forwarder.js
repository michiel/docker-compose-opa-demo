var http = require("http");

const postBody = (msg) => {
  return {
    source: '',
    sourcetype: '',
    index: '',
    fields: (msg.labels || {}),
    event: msg,
  };
};

const postIt = (decision) => {

  var options = {
    hostname: 'splunk',
    port: 8088,
    path: '/services/collector/event',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authentication': '3e6ffd12-0f69-46bb-ad0d-71cffb661a0d',
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

  req.write(JSON.stringify(postBody(decision)));
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
      postIt(body);
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

