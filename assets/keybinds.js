document.addEventListener('keyup', kup)
function kup(e) {
  switch (e.key) {
    case '/':
      document.getElementById('search').focus()
      break;
    default:
      break;
  }
}
