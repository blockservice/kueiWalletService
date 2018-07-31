const XMLHttpRequest = require('xmlhttprequest').XMLHttpRequest;

async function send(url, method, params) {
    const request = {
        method: method,
        params: params,
        id: 11,
        jsonrpc: "2.0"
    };
    return fetchJSON(url, JSON.stringify(request));
}

function fetchJSON(url, payload) {
    return new Promise(function(resolve, reject) {
        const request = new XMLHttpRequest({
            tlsOptions: {
                rejectUnauthorized: false,
                ecdhCurve: 'secp384r1',
            }
        });

        if (payload) {
            request.open('POST', url, true);
            request.setRequestHeader('Content-Type','application/json');
        } else {
            request.open('GET', url, true);
        }

        request.onreadystatechange = function() {
            if (request.readyState !== 4) { return; }
            if (request.status !== 200) {
                return reject(`invalid response - ${request.status}, responseText:${request.responseText}`);
            }
            let json = JSON.parse(request.responseText);
            if(json.error) {
                console.warn(json.error);
                return reject("invalid json");
            }

            resolve(json.result);
        };

        request.onerror = function(error) {
            reject(error);
        };

        try {
            request.send(payload ? payload : null);
        } catch (error) {
            const connectionError = new Error('connection error');
            connectionError['error'] = error;
            reject(connectionError);
        }
    });
}



// 172.31.1.177 18545 18546
// ====================================
// process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = 0;
// send("http://127.0.0.1:8545", 'ews_foo', ["straysh"]).catch(err=>{
//     console.log(err);
// }).then(data=>{
//     console.log(data);
// });

// const WS = require("./ws_provider");
const wait = ms => new Promise(resolve => setTimeout(resolve, ms));
const Web3 = require("web3");
const ws = new Web3.providers.WebsocketProvider("ws://127.0.0.1:8546");

async function run(){
    await wait(500);
    ws.send({
        "jsonrpc": "2.0",
        "id": 101,
        "method": "ews_subscribe",
        "params": ["newToken", "0x3A1CBB00730dDf72f3172C78fb5fbBefFFcc62A7"]
    }, async (err, data)=>{
        console.log(data);
        console.log("-".repeat(64));
        for(let i=0;;i++) {
            let params = {
                "jsonrpc": "2.0",
                "id": 1,
                "method": "ews_subscribe",
                "params": ["pingPong", data.result]
            };
            ws.send(params, (err,data)=>{
                console.log(`ews_subscribe`, data);
            });
            await wait(1000*60*4);
        }
    });

    ws.on("data", function(d){
        console.dir(d, {colors:true});
    });
}
run().catch(err=>{
    console.log(err);
});
