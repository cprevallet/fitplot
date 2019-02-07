// initialize the server

var exec = require('child_process').exec;

// Windows, fitplot with be in the path 
if (navigator.platform.substring(0,3) == "Win") {
	exec('fitplot &', function callback(error, stdout, stderr){
    	if (error) {
        	console.log("child processes failed with error code: " +
            	error.code);
    	}
    	console.log(stdout);
	});
  };
  
// Linux, need to use a relative path
if (navigator.platform.substring(0,3) == "Lin") {
	exec('./fitplot &', function callback(error, stdout, stderr){
    	if (error) {
        	console.log("child processes failed with error code: " +
            	error.code);
    	}
    	console.log(stdout);
	});
  };  
  
// and ...
nw.Window.open('localhost:8080/', {width:1024, height:768}, function(win) {});
