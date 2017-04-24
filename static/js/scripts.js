$( document ).ready(function() {
    getUser()
    getToDoS()
    getNotifys()
});

function getUser() {
	$.ajax({
	  type: "POST",
	  url: "/API/",
	  data: "method=getUser",
	  success: function(msg){
	    var obj = JSON.parse(msg)
	    if (obj.User != undefined) {
	    	$(".not_signed_in").hide()
	    	$(".signed_in").show()
	    	$("#user_name").html(obj.User)	
	    }
	  }
	});	
}

function getToDoS() {
	$.ajax({
		type: "POST",
		url: "/API/",
		data: "method=getToDoS",
		success: function(msg){
			var obj = JSON.parse(msg)

			for (var key in obj) {
				if (obj[key].Status==0) {
					obj[key].Status="inline-block"
				} else {
					obj[key].Status="none"
				}

				$('#todo_template').tmpl(obj[key]).appendTo('#todos');
				$('#form_template').tmpl(obj[key]).appendTo('#modals');
			}
		}
	})
}

function getNotifys() {
	$.ajax({
		type: "POST",
		url: "/API/",
		data: "method=getNotifys",
		success: function(msg){
			var obj = JSON.parse(msg)
			if (msg!="{}") {
				$('#notify_template').tmpl(obj).appendTo('#alerts');
			}
		}
	})	
}

function setDoneToDo(id) {
	$.ajax({
		type: "POST",
		url: "/API/",
		data: "method=editToDo&id="+id+"&status=1",
		success: function(msg){
			location.reload();
		}
	})	
}