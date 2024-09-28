// document.getElementById("user-input-form").addEventListener("submit", function(event) {
//     event.preventDefault();
//
//     const inputText = document.getElementById("user-input").value;
//
//     fetch("/sendInput", {
//         method: "POST",
//         headers: {
//             "Content-Type": "application/json",
//         },
//         body: JSON.stringify({ inputText: inputText })
//     })
//         .then(response => response.text())
//         .then(data => {
//             document.getElementById("response-output").textContent = "Server Response: " + data;
//         })
//         .catch(error => console.error("Error:", error));
// });