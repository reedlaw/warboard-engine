var wsDispatcher = function(url){
    var conn = new WebSocket(url);

    var callbacks = {};

    this.bind = function(event_name, callback){
        callbacks[event_name] = callbacks[event_name] || [];
        callbacks[event_name].push(callback);
        return this;
    };

    this.send = function(event_name, event_data){
        var msg = JSON.stringify({event: event_name, data: event_data});
        conn.send(msg);
        return this;
    };

    conn.onmessage = function(msg){
        try {
            var json = JSON.parse(msg.data);
            dispatch(json.event, json.data);            
        } catch(err) {
            console.log("Error: " + err.message);
        }
    };

    conn.onclose = function(){dispatch('close',null)}
    conn.onopen = function(){dispatch('open',null)}

    var dispatch = function(event_name, message){
        var chain = callbacks[event_name];
        if(typeof chain == 'undefined') return;
        for(var i = 0; i < chain.length; i++){
            chain[i]( message )
        }
    }
}

var ws = new wsDispatcher("ws://localhost:8000/websocket/ws");

// bind to server events
ws.bind('open', function(data){
    console.log("Opened socket");
});

ws.bind('close', function(data){
    console.log("Closed socket");
});

ws.bind('send_message', function(data){
    console.log(data.name + ' says: ' + data.message);
});

ws.bind('time', function(data){
    console.log(data.name + ' time: ' + data.message);
});
