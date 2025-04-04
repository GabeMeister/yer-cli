document.body.addEventListener("htmx:afterSettle", function (evt) {
  new TomSelect("#animal", {
    persist: false,
    createOnBlur: true,
    create: true,
  });
});

document.addEventListener("DOMContentLoaded", function () {
  new TomSelect("#animal", {
    persist: false,
    createOnBlur: true,
    create: true,
  });
});
