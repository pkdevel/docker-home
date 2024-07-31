const lightSwitches = document.querySelectorAll('.light-switch');
if (lightSwitches.length > 0) {
  var dm = localStorage.getItem('dark-mode');
  lightSwitches.forEach((lightSwitch, i) => {
    if (dm === 'true' || (!dm && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      lightSwitch.checked = true;
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
    lightSwitch.addEventListener('change', () => {
      const { checked } = lightSwitch;
      lightSwitches.forEach((el, n) => {
        if (n !== i) {
          el.checked = checked;
        }
      });
      if (lightSwitch.checked) {
        document.documentElement.classList.add('dark');
        localStorage.setItem('dark-mode', true);
      } else {
        document.documentElement.classList.remove('dark');
        localStorage.setItem('dark-mode', false);
      }
    });
  });
}
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
  if (localStorage.getItem('dark-mode')) {
    return;
  }
  if (e.matches) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
});
