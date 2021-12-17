const list = document.querySelector('ul')

var blink_speed = 1000; // every 1000 == 1 second, adjust to suit
var t = setInterval(function () {
  const list = document.querySelector('ul')
  list.style.visibility = (list.style.visibility == 'hidden' ? '' : 'hidden');
}, blink_speed);