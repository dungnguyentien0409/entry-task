(function() {
})()

var token;

function search() {
    var tmp = token
    $.ajax({
            type: "GET",
            dataType: "json",
            url: "http://localhost:49/products?PageIndex=0&PageSize=10&CategoryId=1",
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + token);
            },
            success: function(res) {
                 if (res.error_code == 0) {
                    var arr = JSON.stringify(res.data)

                    AddToResult(arr)
                 }
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                 alert("errot")
            }
        })
}

function AddToResult(data) {
    var table = document.getElementById("tblResult");
    var arr =  JSON.parse(data);
    debugger;

    for(var i = 0; i < arr.length; i++) {
        var item = arr[i];
        debugger;
        var row = table.insertRow(-1);
        var id = row.insertCell(0);
        var categoryId = row.insertCell(1);
        var name = row.insertCell(2);
        var description = row.insertCell(3);
        var price = row.insertCell(4);
        var currency = row.insertCell(5);
        var images = row.insertCell(6);
        var comments = row.insertCell(7);

        id.innerHTML = item.Id;
        categoryId.innerHTML = item.CategoryId;
        name.innerHTML = item.Name;
        description.innerHTML = item.Description;
        price.innerHTML = item.Price;
        currency.innerHTML = item.Currency;
        images.innerHTML = item.Images
        comments.innerHTML = item.Comments;
    }
}

function login() {
    var account = $.trim($("#account").val())
    var password = $.trim($("#password").val())

    if (account == '' || password == '') {
        $("#error").removeClass("hide")
        $("#error").addClass("display")
        return;
    }

    $.ajax({
        type: "POST",
        dataType: "json",
        url: "http://localhost:49/login",
        data: JSON.stringify({
            Account: account,
            Password: password
        }),
        success: function(res) {
             $("message").text(res.message)
             if (res.error_code == 0) {
                token = res.data

                 $("#product").removeClass("hide")
                 $("#product").addClass("display")
             }
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
             $("#error").removeClass("hide")
             $("#error").addClass("display")
        }
    })
}

function reset() {
    $("#product").addClass("hide")
    $("#product").removeClass("display")

    $("#error").addClass("hide")
    $("#error").removeClass("display")

    $("#tblResult").empty()
}

function bin2String(array) {
  var result = "";
  for (var i = 0; i < array.length; i++) {
    result += String.fromCharCode(parseInt(array[i], 2));
  }
  return result;
}