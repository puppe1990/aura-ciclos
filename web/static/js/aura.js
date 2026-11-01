(function () {
  if (window.__auraUi) return;
  window.__auraUi = true;

  function litersLabel(glasses) {
    return (glasses * 0.25).toFixed(1).replace(".", ",") + "L";
  }

  document.addEventListener("click", function (ev) {
    var waterBtn = ev.target.closest("[data-aura-water]");
    if (waterBtn) {
      ev.preventDefault();
      var input = document.getElementById("waterGlasses");
      if (!input) return;
      var n = parseInt(input.value, 10) || 0;
      n += parseInt(waterBtn.getAttribute("data-aura-water"), 10) || 0;
      if (n < 0) n = 0;
      if (n > 20) n = 20;
      input.value = String(n);
      var display = document.getElementById("waterDisplay");
      if (display) {
        display.innerHTML =
          litersLabel(n) +
          ' <span class="text-body-sm text-on-surface-variant font-normal">/ 2,5L</span>';
      }
      return;
    }
    if (ev.target.closest("[data-aura-print]")) {
      ev.preventDefault();
      window.print();
    }
    var timesBtn = ev.target.closest("[data-aura-sex-times]");
    if (timesBtn) {
      ev.preventDefault();
      var input = document.getElementById("sexTimes");
      if (!input) return;
      var n = parseInt(input.value, 10) || 0;
      n += parseInt(timesBtn.getAttribute("data-aura-sex-times"), 10) || 0;
      if (n < 0) n = 0;
      if (n > 20) n = 20;
      input.value = String(n);
      var display = document.getElementById("sexTimesDisplay");
      if (display) display.textContent = String(n);
      return;
    }
    var hideBtn = ev.target.closest("[data-aura-intimate-toggle]");
    if (hideBtn) {
      ev.preventDefault();
      var box = document.getElementById("auraIntimate");
      if (!box) return;
      if (box.hasAttribute("hidden")) box.removeAttribute("hidden");
      else box.setAttribute("hidden", "");
    }
  });
})();
