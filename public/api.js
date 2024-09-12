
const getRequestValue = document.getElementById('req').value;

// Construct Request

function newGetRequest(value, headers = "application/json") {
    return new Request(value, {
        method: "GET",
        headers: {
            "Content-Type": headers,
        },
    });
}

// Async Fetch

// Async Fetch
async function get(request) {
    try {
        const response = await fetch(request);
        const result = await response.json();  // Parse the response as JSON
        console.log("Success:", result);
        document.querySelector('.output-req').textContent = JSON.stringify(result, null, 2);
    } catch (error) {
        console.error("Error:", error);
    }
}

// Button click handler
function handleSubmit() {
    const getRequestValue = document.getElementById('req').value; // Get input value dynamically
    const request = newGetRequest(getRequestValue);               // Construct request
    get(request);                                                 // Make the request
}

// Construct Header