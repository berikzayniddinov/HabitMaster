// main.js

import {
    fetchNotifications,
    addNotification,
    deleteNotification,
    editNotification
} from './notifications.js';
import { currentPage, prevPage, nextPage, updatePaginationControls } from './pagination.js';
import { renderNotifications } from './utils.js';

const notificationContainer = document.getElementById('notifications-container');

async function getNotifications() {
    const filter = document.getElementById('notification-filter')?.value || '';
    const sort = document.getElementById('notification-sort')?.value || '';

    try {
        const notifications = await fetchNotifications(filter, sort, currentPage);
        renderNotifications(notificationContainer, notifications, editNotification, deleteNotification);
        updatePaginationControls();
    } catch (error) {
        console.error(error.message);
        alert('Failed to fetch notifications.');
    }
}

document.querySelector('.add-button').addEventListener('click', async () => {
    const message = prompt('Enter notification message:');
    const scheduledAt = prompt('Enter scheduled date and time (YYYY-MM-DDTHH:mm:ss):');
    if (!message || !scheduledAt) {
        alert('Both message and scheduled date are required!');
        return;
    }

    try {
        await addNotification({ message, scheduled_at: scheduledAt });
        alert('Notification successfully added!');
        getNotifications();
    } catch (error) {
        console.error(error.message);
        alert('Failed to add notification.');
    }
});

window.prevPage = () => prevPage(getNotifications);
window.nextPage = () => nextPage(getNotifications);

window.onload = getNotifications;
