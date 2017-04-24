$( document ).ready(function() {
    console.log( "ready!" );
    getUser()
});

function getUser() {
	$.ajax({
	  type: "POST",
	  url: "/API/",
	  data: "method=getUser",
	  success: function(msg){
	    var obj = JSON.parse(msg)
	    console.log(obj.User)
	    if (obj.User != undefined) {
	    	$(".not_signed_in").hide()
	    	$(".signed_in").show()
	    	$("#user_name").html(obj.User)	
	    }
	  }
	});	
}