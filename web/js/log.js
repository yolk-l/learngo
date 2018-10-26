function do_print() {
    $('#log_content').append("do 'update server' \n");
    var begin = 1;
    var offset = 10;
    var timer = setInterval(function() {
        var url = 'http://'+location.host+'/ajax?begin='+begin+'&offset='+offset;
        $.ajax({
            url: url,
            type: 'GET',
            success: function(ret) {
                if(ret.length > 0 ) {
                    $("#log_content").append(ret);
                    for(i=0;i<ret.length;i++) {
                        if(ret[i] == "\n") {
                            begin++;
                        }
                    }
                    $("#log").scrollTop($('#log').scrollTop()+300);
                }
                else {
                    clearInterval(timer)
                }

            },
            error: function(ret) {
                $('#error').append('update failed! <br />' + ret);
                clearInterval(timer);
                return false;
            }
        });
    }, 500); 
}
function do_update() {
    var result;
    $('#log_content').append("do 'svn update' \n");
    $.ajax({
        url: 'http://'+location.host+'/update',
        type: 'GET',
        success: function(ret) {
            console.log("...", ret);
            if(ret != undefined) {
                $('#log_content').append(ret);
                do_print();
            }
        },
        error: function(ret) {
            $('#error').append('update failed! <br />' + ret);
        }
    });
}

function onButton1() {
    do_update();
}

$(document).ready(function() {
    var oBnt = $('#bnt1');
    oBnt.click(onButton1);
});
