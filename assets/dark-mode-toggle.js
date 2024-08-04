const ls = document.querySelectorAll('.light-switch');
const cfg_dm = localStorage.getItem('dark-mode');
const wm = window.matchMedia('(prefers-color-scheme: dark)');
const is_dark = cfg_dm === 'true' || (!cfg_dm && wm.matches);
const r = document.documentElement.classList;
if (is_dark) {
  r.add('dark');
} else {
  r.remove('dark');
}
if (ls.length > 0) {
  ls.forEach((lightSwitch, i) => {
    if (is_dark) {
      lightSwitch.checked = true;
    }
    lightSwitch.addEventListener('change', () => {
      const { checked } = lightSwitch;
      ls.forEach((el, n) => {
        if (n !== i) {
          el.checked = checked;
        }
      });
      if (lightSwitch.checked) {
        r.add('dark');
        localStorage.setItem('dark-mode', true);
      } else {
        r.remove('dark');
        localStorage.setItem('dark-mode', false);
      }
    });
  });
}
wm.addEventListener('change', (e) => {
  if (localStorage.getItem('dark-mode')) {
    return;
  }
  if (e.matches) {
    r.add('dark');
  } else {
    r.remove('dark');
  }
});
