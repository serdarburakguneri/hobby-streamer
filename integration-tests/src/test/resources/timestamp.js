function fn() {
  var generateTimestamp = function() {
    var now = new Date();
    var pad = function(n){ return n < 10 ? '0' + n : n; };
    return now.getFullYear().toString() +
      pad(now.getMonth() + 1) +
      pad(now.getDate()) +
      pad(now.getHours()) +
      pad(now.getMinutes()) +
      pad(now.getSeconds());
  };
  return { generateTimestamp: generateTimestamp };
}