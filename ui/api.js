function updateFormAction() {
    // Get the value from the input field
    const inputField = document.getElementById("req").value;

    // Select the form element
    const form = document.querySelector(".api-tool");

    // Update the hx-get attribute of the form
    form.setAttribute("hx-get", inputField);

    // Optionally log to console for debugging
    console.log("Updated hx-get to:", inputField);
}