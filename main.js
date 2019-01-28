// initialize the server
var exec = require('child_process').exec;
exec('./fitplot &', function callback(error, stdout, stderr){
    // result
});


// and ...
nw.Window.open('localhost:8080/', {width:1024, height:768}, function(win) {});
