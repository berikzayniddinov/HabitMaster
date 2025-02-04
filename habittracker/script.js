// Создание карточки элемента
function createCard(id, title, description, date, type, extra = '', actions = []) {
    const container = document.getElementById(`${type}-container`);
    const card = document.createElement('div');
    card.classList.add('card');
    card.dataset.name = title || 'No Name'; // Привязка к текущему названию

    card.innerHTML = `
        <h3>${title || 'No Title'}</h3>
        <p>${description || 'No Description'}</p>
        <p><strong>Date:</strong> ${date || 'No Date'}</p>
        ${extra ? `<p>${extra}</p>` : ''}
        <div class="card-actions"></div>
    `;

    const actionsContainer = card.querySelector('.card-actions');
    actions.forEach(action => actionsContainer.appendChild(action));
    container.appendChild(card);
}


// Создание кнопки действия
function createActionButton(text, className, onclick) {
    const button = document.createElement('button');
    button.textContent = text;
    button.classList.add(className);
    button.onclick = onclick;
    return button;
}

document.getElementById("menuToggle").addEventListener("click", () => {
    const sideMenu = document.getElementById("sideMenu");
    sideMenu.style.display = sideMenu.style.display === "block" ? "none" : "block";
});
document.addEventListener("DOMContentLoaded", function () {
    const menuToggle = document.getElementById("menuToggle");
    const sideMenu = document.getElementById("sideMenu");

    if (!menuToggle || !sideMenu) {
        console.error("Menu elements not found!");
        return;
    }

    menuToggle.addEventListener("click", function () {
        sideMenu.classList.toggle("active"); // Добавляет/убирает класс
    });
});

