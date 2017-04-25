$( document ).ready(function() {
    getUser();
    getToDoS();
    getNotifys();
    attachHandlers()
});

function getUser() {
	$.ajax({
	  type: "POST",
	  url: "/API/",
	  data: "method=getUser",
	  success: function(msg){
          var obj = JSON.parse(msg);
	    if (obj.User != undefined) {
            $(".not_signed_in").hide();
            $(".signed_in").show();
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
            var obj = JSON.parse(msg);
            var count = 0;
			for (var key in obj) {
                count++;
				if (obj[key].Status==0) {
					obj[key].Status="inline-block"
				} else {
					obj[key].Status="none"
				}

				$('#todo_template').tmpl(obj[key]).appendTo('#todos');
				$('#form_template').tmpl(obj[key]).appendTo('#modals');
			}
            if (count == 0 && msg != "{}") {
                var notify = {
                    Title: "Привет!",
                    Text: "Для добавления нового ToDo: в верхнем меню нажмите \"Добавить ToDo:\"",
                    Status: "success"
                };
                $('#notify_template').tmpl(notify).appendTo('#alerts');
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
            var obj = JSON.parse(msg);
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

function attachHandlers() {
    $("#login_submit").click(function () {
        login = $("#login_login").val();
        pass = $("#login_pass").val();
        $.ajax({
            type: "POST",
            url: "/login/",
            data: "login=" + login + "&pass=" + pass,
            success: function (msg) {
                if (msg != "") {
                    $('#login_alert').html(msg);
                    $('#login_alert').slideDown(100)
                } else {
                    location.reload();
                }
            }
        })
    });

    $("#register_submit").click(function () {
        login = $("#register_login").val();
        pass = $("#register_pass").val();
        $.ajax({
            type: "POST",
            url: "/register/",
            data: "login=" + login + "&pass=" + pass,
            success: function (msg) {
                if (msg != "") {
                    $('#register_alert').html(msg);
                    $('#register_alert').slideDown(100)
                } else {
                    location.reload();
                }
            }
        })
    })

}