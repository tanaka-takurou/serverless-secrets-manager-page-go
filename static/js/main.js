$(document).ready(function() {
  GetSecret();
});

var GetSecret = function() {
  const data = {action: "getsecret"};
  request(data, (res)=>{
    console.log(res);
    if (!!res && !!res.message && res.message.length > 0) {
      secretData = JSON.parse(res.message);
    }
    $("#result").text(JSON.stringify(secretData, null, "\t"));
    $("#info").removeClass("hidden").addClass("visible");
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  });
};

var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};
var App = { url: location.origin + {{ .ApiPath }} };
