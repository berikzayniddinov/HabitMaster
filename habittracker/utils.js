// utils.js

export function createCard(notification, editCallback, deleteCallback) {
    const card = document.createElement('div');
    card.className = 'card';
    card.innerHTML = `
        <h3>${notification.message}</h3>
        <p><strong>Scheduled At:</strong> ${new Date(notification.scheduled_at).toLocaleString()}</p>
        <p><strong>Is Sent:</strong> ${notification.is_sent ? 'Yes' : 'No'}</p>
        <div class="card-actions">
            <button class="edit" onclick="${editCallback}(${notification.id})">Edit</button>
            <button class="delete" onclick="${deleteCallback}(${notification.id})">Delete</button>
        </div>
    `;
    return card;
}

export function renderNotifications(container, notifications, editCallback, deleteCallback) {
    container.innerHTML = '';
    notifications.forEach(notification => {
        const card = createCard(notification, editCallback, deleteCallback);
        container.appendChild(card);
    });
}
