// pagination.js

export let currentPage = 1;

export function prevPage(callback) {
    if (currentPage > 1) {
        currentPage--;
        callback();
    }
}

export function nextPage(callback) {
    currentPage++;
    callback();
}

export function updatePaginationControls() {
    const paginationContainer = document.getElementById('notifications-pagination');
    paginationContainer.innerHTML = `
        <button onclick="prevPage()">Previous</button>
        <button onclick="nextPage()">Next</button>
    `;
}
