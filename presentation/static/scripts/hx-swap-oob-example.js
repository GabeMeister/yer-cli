// console.log("\n\n***** Script begin *****\n", "\n\n");

// document.body.addEventListener("htmx:afterSwap", function (evt) {
//   console.log("\n\n***** afterSwap *****\n", window.animalInput, "\n\n");
//   const animalInput = document.querySelector("#animal");
//   console.log("\n\n***** animalInput *****\n", animalInput, "\n\n");
//   window.animalInput = new TomSelect(animalInput, {
//     persist: false,
//     createOnBlur: true,
//     create: true,
//   });
// });

// document.body.addEventListener("htmx:beforeSwap", function (evt) {
//   console.log(
//     "\n\n***** window.animalInput *****\n",
//     window.animalInput,
//     "\n\n"
//   );
//   window.animalInput.destroy();
//   window.animalInput = null;
// });

// document.addEventListener("DOMContentLoaded", function () {
//   const animalInput = document.querySelector("#animal");
//   window.animalInput = new TomSelect(animalInput, {
//     persist: false,
//     createOnBlur: true,
//     create: true,
//   });
// });
console.log("\n\n***** Script begin *****\n", "\n\n");

document.body.addEventListener("htmx:beforeSwap", function (evt) {
  console.log("\n\n***** beforeSwap *****\n", window.animalInput, "\n\n");

  // Check if TomSelect instance exists before destroying it
  if (window.animalInput) {
    window.animalInput.destroy();
    window.animalInput = null; // Clear reference
  }
});

document.body.addEventListener("htmx:afterSwap", function (evt) {
  console.log("\n\n***** afterSwap *****\n");

  const animalElem = document.querySelector("#animal");
  console.log("\n\n***** animalElem *****\n", animalElem, "\n\n");

  // Ensure the new input element exists before reinitializing TomSelect
  if (animalElem) {
    window.animalInput = new TomSelect("#animal", {
      persist: false,
      createOnBlur: true,
      create: true,
    });
    const animalElem = document.querySelector("#animal");
    console.log("\n\n***** THIS IS IT *****\n", animalElem, "\n\n");
    animalElem.classList.add("tomselected", "ts-hidden-accessible");
    console.log("\n\n***** THIS IS IT 2 *****\n", animalElem, "\n\n");
  }
});

// Initialize TomSelect on page load
document.addEventListener("DOMContentLoaded", function () {
  const animalElem = document.querySelector("#animal");
  if (animalElem) {
    window.animalInput = new TomSelect("#animal", {
      persist: false,
      createOnBlur: true,
      create: true,
    });
  }
});
