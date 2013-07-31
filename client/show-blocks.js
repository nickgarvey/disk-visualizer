var TOTAL_SECTORS = 1953525168;
var NUM_ROWS = 10;
var SPACE_BETWEEN_LINES = 10;

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
        fadeTo(.25);

        window.setTimeout(function() {
            fadeTo(.1);

            window.setTimeout(function() {
                fadeTo(.05);

                window.setTimeout(function() {
                    fadeTo(.01);

                    window.setTimeout(function() {
                        fadeTo(0);

                    }, 1000);
                }, 1000);
            }, 1000);
        }, 1000);
    }, 1000);

};

window.onload = function() {
    var canvas = document.getElementById('sda-div');
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    var lineHeight = Math.ceil(canvas.height / NUM_ROWS - SPACE_BETWEEN_LINES);

    var ctx = canvas.getContext("2d");

    for (var rowNum = 0; rowNum < NUM_ROWS; rowNum++) {
        var y = rowNum * (lineHeight + SPACE_BETWEEN_LINES);
        ctx.fillRect(0, y, canvas.width, 1);

        y = rowNum * (lineHeight + SPACE_BETWEEN_LINES) + lineHeight + 2;
        ctx.fillRect(0, y, canvas.width, 1);
    }

    socket.onmessage = function(msg) {
        var trace = JSON.parse(msg.data);

        var percentInDisk = trace.Sector / TOTAL_SECTORS;
        var rowNum = Math.floor(NUM_ROWS * percentInDisk);

        var x = Math.floor(canvas.width * (NUM_ROWS * percentInDisk - rowNum));
        var y = rowNum * (lineHeight + SPACE_BETWEEN_LINES) + 1
        var width = Math.max(2, canvas.width * trace.Blocks / TOTAL_SECTORS);

        var color = new rgba(0, 0, 0, 1);
        if (trace.IoType == "W" || trace.IoType == "WS") {
            color = new rgba(255, 0, 0, 1);
        } else if (trace.IoType == "R" || trace.IoType == "RM") {
            color = new rgba(0, 255, 0, 1);
        } else {
            console.log("Unexpected type: " + trace.IoType);
            return;
        }

        ctx.fillStyle = color.toString();

        ctx.fillRect(x, y, width, lineHeight);
        
        clearRect(ctx, x, y, width, lineHeight, color);
    };
};
