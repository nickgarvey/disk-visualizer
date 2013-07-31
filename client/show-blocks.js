var TOTAL_SECTORS = 1953525168;
var NUM_ROWS = 5;

function rgba(r, g, b, a) {
    this.r = r;
    this.g = g;
    this.b = b;
    this.a = a;
};

rgba.prototype.toString = function rgbaToString() {
    return "rgba(" +
        Math.floor(this.r) + ", " +
        Math.floor(this.g) + ", " +
        Math.floor(this.b) + ", " +
        this.a.toFixed(1) + ")";
};

var socket = new WebSocket("ws://ngarveySiris:8080/ws");
socket.onerror = function(msg) {
    console.log("problem");
    console.log(msg.data);
};

var clearRect = function(ctx, x, y, width, height, color) {

    var fadeTo = function(a) {
        ctx.fillStyle = "#FFF";
        ctx.fillRect(x, y, width, height);

        color.a = a;
        ctx.fillStyle = color.toString();
        ctx.fillRect(x, y, width, height);
    };

    window.setTimeout(function() {
        fadeTo(.5);

        window.setTimeout(function() {
            fadeTo(.25);

            window.setTimeout(function() {
                fadeTo(0);

            }, 2000);
        }, 2000);
    }, 2000);

};

window.onload = function() {
    var canvas = document.getElementById('sda-div');
    canvas.width = window.innerWidth;
    var canvasWidth = canvas.width;
    var ctx = canvas.getContext("2d");
    socket.onmessage = function(msg) {
        var trace = JSON.parse(msg.data);

        var percentInDisk = trace.Sector / TOTAL_SECTORS;
        var rowNum = Math.floor(NUM_ROWS * percentInDisk);

        var x = Math.floor(canvasWidth * (NUM_ROWS * percentInDisk - rowNum));
        var y = rowNum * 150;
        var width = Math.max(2, canvasWidth * trace.Blocks / TOTAL_SECTORS);
        var height = 100;

        var color = new rgba(0, 0, 0, 1);
        if (trace.IoType == "W" || trace.IoType == "WS") {
            color = new rgba(255, 0, 0, 1);
        } else if (trace.IoType == "R" || trace.IoType == "RM") {
            color = new rgba(0, 255, 0, 1);
        } else {
            console.log("Unexpected type: " + trace.IoType);
        }

        ctx.fillStyle = color.toString();

        ctx.fillRect(x, y, width, height);
        
        clearRect(ctx, x, y, width, height, color);
    };
};
