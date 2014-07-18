var lastReceived = 0;
var isWait = false;

var fetch = function () {
    if (isWait) return;
    isWait = true;
    $.getJSON("/lp/fetch?lastReceived=" + lastReceived+"&roomid=" + $('#roomid').text(), function (data) {
        if (data != null) {
        $.each(data, function (i, event) {
            switch (event.Type) {
            case 0: // JOIN
                if (event.User.Name == $('#uname').text()) {
                    $("#chatbox li").first().before("<li>You joined the chat room.</li>");
                } else {
                    $("#chatbox li").first().before("<li>" + event.User.Name + " joined the chat room.</li>");
                }
                break;
            case 1: // LEAVE
                $("#chatbox li").first().before("<li>" + event.User.Name + " left the chat room.</li>");
                break;
            case 2: // MESSAGE
                $("#chatbox li").first().before("<li><b>" + event.User.Name + "</b>: " + event.Content + "</li>");
                break;
            }

            lastReceived = event.Timestamp;
        });
        }
        isWait = false;
    },function(err){
        isWait=false;
    });
}

// Call fetch every 3 seconds
setInterval(fetch, 3000);
 

$(document).ready(function () {
    
    var postConecnt = function () {
        var uname = $('#uname').text();
        var roomid= $('#roomid').text();
        var content = $('#sendbox').val();
        $.post("/lp/post", {
            uname: uname,
            content: content,
            roomid:roomid
        });
        $('#sendbox').val("");
    }

    $('#sendbtn').click(function () {
        postConecnt();
    });  
   fetch();
});